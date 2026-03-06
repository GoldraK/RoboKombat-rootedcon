package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"gopkg.in/yaml.v3"
)

const (
	// XInput button bitmasks
	XINPUT_GAMEPAD_LEFT_SHOULDER  = 0x0100 // LB/L1
	XINPUT_GAMEPAD_RIGHT_SHOULDER = 0x0200 // RB/R1
	// Triggers are analog: LeftTrigger, RightTrigger (0-255)
	TRIGGER_THRESHOLD = 128 // Trigger pressed if > this value
)

// Configuration structures
type Config struct {
	InputMode string `yaml:"inputMode"` // "xbox" or "keyboard"
	Motors    struct {
		InvertA bool `yaml:"invertA"`
		InvertB bool `yaml:"invertB"`
	} `yaml:"motors"`
	Servos struct {
		Servo1 ServoRange `yaml:"servo1"`
		Servo2 ServoRange `yaml:"servo2"`
	} `yaml:"servos"`
	Keyboard KeyboardControls `yaml:"keyboard"`
}

type KeyboardControls struct {
	MotorA struct {
		Forward  string `yaml:"forward"`
		Backward string `yaml:"backward"`
	} `yaml:"motorA"`
	MotorB struct {
		Forward  string `yaml:"forward"`
		Backward string `yaml:"backward"`
	} `yaml:"motorB"`
	Servo1 struct {
		Min string `yaml:"min"`
		Max string `yaml:"max"`
	} `yaml:"servo1"`
	Servo2 struct {
		Min string `yaml:"min"`
		Max string `yaml:"max"`
	} `yaml:"servo2"`
	Stop string `yaml:"stop"`
}

type ServoRange struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

func loadConfig(filename string) (*Config, error) {
	config := &Config{}

	// Set defaults
	config.InputMode = "xbox"
	config.Servos.Servo1.Min = 0
	config.Servos.Servo1.Max = 180
	config.Servos.Servo2.Min = 0
	config.Servos.Servo2.Max = 180

	// Keyboard defaults
	config.Keyboard.MotorA.Forward = "w"
	config.Keyboard.MotorA.Backward = "s"
	config.Keyboard.MotorB.Forward = "i"
	config.Keyboard.MotorB.Backward = "k"
	config.Keyboard.Servo1.Min = "a"
	config.Keyboard.Servo1.Max = "d"
	config.Keyboard.Servo2.Min = "j"
	config.Keyboard.Servo2.Max = "l"
	config.Keyboard.Stop = "space"

	data, err := os.ReadFile(filename)
	if err != nil {
		// If config file doesn't exist, return defaults
		if os.IsNotExist(err) {
			log.Printf("Config file %s not found, using defaults", filename)
			return config, nil
		}
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return config, nil
}

func main() {
	// CLI arguments
	apiHost := flag.String("host", "192.168.4.1", "API server IP address")
	controllerID := flag.Int("controller", 0, "Xbox controller ID (0-3)")
	debug := flag.Bool("debug", false, "Enable debug logs for buttons and axes")
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *apiHost == "" {
		log.Fatal("API host IP is required. Use -host <IP>")
	}

	fmt.Printf("Connecting to RoboCombat at %s\n", *apiHost)
	fmt.Printf("Configuration: Input=%s, Motor A inverted=%v, Motor B inverted=%v\n",
		config.InputMode, config.Motors.InvertA, config.Motors.InvertB)

	httpClient := &http.Client{Timeout: 500 * time.Millisecond}

	// Channel for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Extra debugging to check why it's exiting
	go func() {
		s := <-stop
		fmt.Printf("\nReceived signal: %v. Stopping efforts...\n", s)
		sendStop(*apiHost, httpClient)
		os.Exit(0)
	}()

	// Start appropriate input mode
	if config.InputMode == "keyboard" {
		runKeyboardMode(*apiHost, config, *debug, httpClient)
	} else {
		runXboxMode(*apiHost, config, *controllerID, *debug, httpClient)
	}
}

func runXboxMode(apiHost string, config *Config, controllerID int, debug bool, httpClient *http.Client) {
	// Test controller connection
	_, err := XInputGetState(uint32(controllerID))
	if err != nil {
		log.Fatalf("Failed to connect to controller %d: %v\nMake sure your Xbox controller is connected.", controllerID, err)
	}

	fmt.Printf("Xbox controller %d connected successfully\n", controllerID)
	fmt.Println("Tank mode: Left stick = Motor A, Right stick = Motor B")
	fmt.Println("Servo control: L1/L2 = Servo 1, R1/R2 = Servo 2")

	if debug {
		log.Println("DEBUG mode enabled - will show all controller activity")
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	lastA, lastB := 0, 0
	lastServo1, lastServo2 := 90, 90 // Start at middle position
	loopCount := 0
	var lastButtons uint16 = 0

	fmt.Println("Starting main loop... Press Ctrl+C to exit")

	for {
		<-ticker.C
		loopCount++
		if debug && loopCount%20 == 0 {
			log.Printf("Heartbeat: %d ticks", loopCount)
		}

		state, err := XInputGetState(uint32(controllerID))
		if err != nil {
			if debug {
				log.Printf("Error reading controller: %v", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}

		if debug {
			// Button monitoring
			var pressed []int
			for i := 0; i < 16; i++ {
				if (state.Gamepad.Buttons>>uint(i))&1 != 0 {
					pressed = append(pressed, i)
				}
			}

			// Detect button changes
			if state.Gamepad.Buttons != lastButtons || len(pressed) > 0 {
				fmt.Printf("\rDEBUG - Buttons: %v | LX:%d LY:%d RX:%d RY:%d LT:%d RT:%d          \n",
					pressed,
					state.Gamepad.ThumbLX,
					state.Gamepad.ThumbLY,
					state.Gamepad.ThumbRX,
					state.Gamepad.ThumbRY,
					state.Gamepad.LeftTrigger,
					state.Gamepad.RightTrigger,
				)
			}
			lastButtons = state.Gamepad.Buttons
		}

		// Tank mode: Left stick controls motor A, right stick controls motor B
		// Apply deadzone to each stick
		leftY := applyDeadzone(state.Gamepad.ThumbLY, 7849)
		rightY := applyDeadzone(state.Gamepad.ThumbRY, 8689)

		// Scale from -32768..32767 to -255..255
		// XInput: Up=positive, Down=negative
		a := scaleAxis(leftY)
		b := scaleAxis(rightY)

		// Apply motor inversions from config
		if config.Motors.InvertA {
			a = -a
		}
		if config.Motors.InvertB {
			b = -b
		}

		// Servo control with shoulder buttons and triggers
		servo1 := lastServo1
		servo2 := lastServo2

		// L1 (left shoulder) -> Servo 1 to min position
		// L2 (left trigger) -> Servo 1 to max position
		if state.Gamepad.Buttons&XINPUT_GAMEPAD_LEFT_SHOULDER != 0 {
			servo1 = config.Servos.Servo1.Min
		} else if state.Gamepad.LeftTrigger > TRIGGER_THRESHOLD {
			servo1 = config.Servos.Servo1.Max
		}

		// R1 (right shoulder) -> Servo 2 to min position
		// R2 (right trigger) -> Servo 2 to max position
		if state.Gamepad.Buttons&XINPUT_GAMEPAD_RIGHT_SHOULDER != 0 {
			servo2 = config.Servos.Servo2.Min
		} else if state.Gamepad.RightTrigger > TRIGGER_THRESHOLD {
			servo2 = config.Servos.Servo2.Max
		}

		// Send control if anything changed
		if a != lastA || b != lastB || servo1 != lastServo1 || servo2 != lastServo2 {
			if debug {
				log.Printf("Sending: Motors A=%d B=%d | Servos S1=%d S2=%d", a, b, servo1, servo2)
			}
			go sendControlWithServos(apiHost, httpClient, a, b, servo1, servo2)
			lastA, lastB = a, b
			lastServo1, lastServo2 = servo1, servo2
		}
	}
}

func runKeyboardMode(apiHost string, config *Config, debug bool, httpClient *http.Client) {
	fmt.Println("Keyboard mode enabled")
	fmt.Printf("Controls: %s/%s=Motor A  %s/%s=Motor B  %s/%s=Servo1  %s/%s=Servo2  %s=Stop\n",
		config.Keyboard.MotorA.Forward, config.Keyboard.MotorA.Backward,
		config.Keyboard.MotorB.Forward, config.Keyboard.MotorB.Backward,
		config.Keyboard.Servo1.Min, config.Keyboard.Servo1.Max,
		config.Keyboard.Servo2.Min, config.Keyboard.Servo2.Max,
		config.Keyboard.Stop)

	if err := keyboard.Open(); err != nil {
		log.Fatalf("Failed to open keyboard: %v", err)
	}
	defer keyboard.Close()

	if debug {
		log.Println("DEBUG mode enabled - will show all key presses")
	}

	// State tracking
	keysPressed := make(map[string]bool)
	lastA, lastB := 0, 0
	lastServo1, lastServo2 := 90, 90

	// Motor speed increment
	const motorSpeed = 220

	fmt.Println("Starting keyboard control... Press Ctrl+C to exit")

	// Read keyboard in separate goroutine
	keyEvents := make(chan keyboard.KeyEvent, 10)
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				return
			}
			keyEvents <- keyboard.KeyEvent{Rune: char, Key: key}
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case event := <-keyEvents:
			keyStr := keyEventToString(event)

			if debug {
				log.Printf("Key pressed: %s", keyStr)
			}

			// Toggle key state
			if keyStr != "" {
				keysPressed[keyStr] = true
			}

			// Stop key stops everything
			if keyStr == config.Keyboard.Stop {
				keysPressed = make(map[string]bool)
				sendStop(apiHost, httpClient)
				lastA, lastB = 0, 0
				lastServo1, lastServo2 = 90, 90
				if debug {
					log.Println("STOP - All motors stopped")
				}
			}

		case <-ticker.C:
			// Calculate motor speeds based on keys pressed
			a, b := 0, 0

			// Motor A
			if keysPressed[config.Keyboard.MotorA.Forward] {
				a = motorSpeed
			} else if keysPressed[config.Keyboard.MotorA.Backward] {
				a = -motorSpeed
			}

			// Motor B
			if keysPressed[config.Keyboard.MotorB.Forward] {
				b = motorSpeed
			} else if keysPressed[config.Keyboard.MotorB.Backward] {
				b = -motorSpeed
			}

			// Apply motor inversions from config
			if config.Motors.InvertA {
				a = -a
			}
			if config.Motors.InvertB {
				b = -b
			}

			// Servo control - return to center (90) when no key pressed
			servo1 := 90
			servo2 := 90

			// Servo 1
			if keysPressed[config.Keyboard.Servo1.Min] {
				servo1 = config.Servos.Servo1.Min
			} else if keysPressed[config.Keyboard.Servo1.Max] {
				servo1 = config.Servos.Servo1.Max
			}

			// Servo 2
			if keysPressed[config.Keyboard.Servo2.Min] {
				servo2 = config.Servos.Servo2.Min
			} else if keysPressed[config.Keyboard.Servo2.Max] {
				servo2 = config.Servos.Servo2.Max
			}

			// Send control if anything changed
			if a != lastA || b != lastB || servo1 != lastServo1 || servo2 != lastServo2 {
				if debug {
					log.Printf("Sending: Motors A=%d B=%d | Servos S1=%d S2=%d", a, b, servo1, servo2)
				}
				go sendControlWithServos(apiHost, httpClient, a, b, servo1, servo2)
				lastA, lastB = a, b
				lastServo1, lastServo2 = servo1, servo2
			}

			// Clear key presses after processing (motors stop when key released)
			keysPressed = make(map[string]bool)
		}
	}
}

func scaleAxis(value int16) int {
	// Scale from -32768..32767 to -255..255
	// Use full range for maximum precision
	if value > 0 {
		return int(value) * 255 / 32767
	} else {
		return int(value) * 255 / 32768
	}
}

func sendControl(host string, client *http.Client, a, b int) {
	u := fmt.Sprintf("http://%s/api/motors?a=%d&b=%d", host, a, b)
	resp, err := client.Get(u)
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return
	}
	resp.Body.Close()
}

func sendControlWithServos(host string, client *http.Client, a, b, s1, s2 int) {
	u := fmt.Sprintf("http://%s/api/control?a=%d&b=%d&s1=%d&s2=%d", host, a, b, s1, s2)
	resp, err := client.Get(u)
	if err != nil {
		log.Printf("Error sending control command: %v", err)
		return
	}
	resp.Body.Close()
}

func sendStop(host string, client *http.Client) {
	u := fmt.Sprintf("http://%s/api/stop", host)
	resp, err := client.Get(u)
	if err == nil {
		resp.Body.Close()
	}
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

func keyEventToString(event keyboard.KeyEvent) string {
	// Handle special keys
	switch event.Key {
	case keyboard.KeySpace:
		return "space"
	case keyboard.KeyEnter:
		return "enter"
	case keyboard.KeyTab:
		return "tab"
	case keyboard.KeyEsc:
		return "esc"
	case keyboard.KeyArrowUp:
		return "up"
	case keyboard.KeyArrowDown:
		return "down"
	case keyboard.KeyArrowLeft:
		return "left"
	case keyboard.KeyArrowRight:
		return "right"
	}

	// Handle regular characters
	if event.Rune != 0 {
		char := event.Rune
		if char >= 'A' && char <= 'Z' {
			char = char + ('a' - 'A')
		}
		return string(char)
	}

	return ""
}
