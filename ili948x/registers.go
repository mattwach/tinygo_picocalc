package ili948x

const ( // ILI9488 Datasheet, pp. 140-147
	CMD_NOP     uint8 = 0x00 // No Operation
	CMD_SWRESET       = 0x01 // Software Reset

	CMD_RDDIDIF    = 0x04 // Read Display Identification Information
	CMD_RDNUMED    = 0x05 // Read Number of the Errors on DSI
	CMD_RDDST      = 0x09 // Read Display Status
	CMD_RDDPM      = 0x0a // Read Display Power Mode
	CMD_RDDMADCTRL = 0x0b // Read Display MADCTRL
	CMD_RDDCOLMOD  = 0x0c // Read Display Pixel Format
	CMD_RDDIM      = 0x0d // Read Display Image Mode
	CMD_RDDSM      = 0x0e // Read Display Signal Mode
	CMD_RDDSDR     = 0x0f // Read Display Self-Diagnostic Result

	CMD_SLPIN   = 0x10 // Enter Sleep Mode
	CMD_SLPOUT  = 0x11 // Sleep Out
	CMD_PTLON   = 0x12 // Partial Mode ON
	CMD_NORON   = 0x13 // Normal Display Mode ON
	CMD_INVOFF  = 0x20 // Display Inversion OFF
	CMD_INVON   = 0x21 // Display Inversion ON
	CMD_ALLPOFF = 0x22 // All Pixels OFF
	CMD_ALLPON  = 0x23 // All Pixels ON
	CMD_DISOFF  = 0x28 // Display OFF
	CMD_DISON   = 0x29 // Display ON

	CMD_CASET    = 0x2a // Column Address Set
	CMD_PASET    = 0x2b // Page Address Set
	CMD_RAMWR    = 0x2c // Memory Write
	CMD_RAMRD    = 0x2e // Memory Read
	CMD_PLTAR    = 0x30 // Partial Area
	CMD_VSCRDEF  = 0x33 // Vertical Scrolling Definition
	CMD_TEOFF    = 0x34 // Tearing Effect Line OFF
	CMD_TEON     = 0x35 // Tearing Effect Line ON
	CMD_MADCTRL  = 0x36 // Memory Access Control
	CMD_VSCRSADD = 0x37 // Vertical Scrolling Start Address
	CMD_IDMOFF   = 0x38 // Idle Mode OFF
	CMD_IDMON    = 0x39 // Idle Mode ON
	CMD_PIXFMT   = 0x3a // COLMOD: Interface Pixel Format
	CMD_RAMWRC   = 0x3c // Memory Write Continue
	CMD_RAMRDRC  = 0x3e // Memory Read Continue
	CMD_TESLWR   = 0x44 // Write Tear Scan Line
	CMD_TESLRD   = 0x45 // Read Tear Scan Line
	CMD_WRDISBV  = 0x51 // Write Display Brightness Value
	CMD_RDDISBV  = 0x52 // Read Display Brightness Value
	CMD_WRCTRLD  = 0x53 // Write CTRL Display Value
	CMD_RDCTRLD  = 0x54 // Read CTRL Display Value
	CMD_WRCABC   = 0x55 // Write Content Adaptive Brightness Control Value
	CMD_RDCABC   = 0x56 // Read Content Adaptive Brightness Control Value
	CMD_WRCABCMB = 0x5e // Write CABC Minimum Brightness
	CMD_RDCABCMB = 0x5f // Read CABC Minimum Brightness
	CMD_RDABCSDR = 0x68 // Read Automatic Brightness Control Self-diagnostic Result

	CMD_IFMODE   = 0xb0 // Interface Mode Control
	CMD_FRMCTRL1 = 0xb1 // Frame Rate Control (In Normal Mode/Full Colors)
	CMD_FRMCTRL2 = 0xb2 // Frame Rate Control (In Idle Mode/8 colors)
	CMD_FRMCTRL3 = 0xb3 // Frame Rate control (In Partial Mode/Full Colors)
	CMD_INVCTRL  = 0xb4 // Display Inversion Control
	CMD_PRCTRL   = 0xb5 // Blanking Porch Control
	CMD_DISCTRL  = 0xb6 // Display Function Control
	CMD_ETMOD    = 0xb7 // Entry Mode Set
	CMD_CECTRL1  = 0xb9 // Color Enhancement Control 1
	CMD_CECTRL2  = 0xba // Color Enhancement Control 2
	CMD_HSLCTRL  = 0xbe // HS Lanes Control

	CMD_PWCTRL1   = 0xc0 // Power Control 1
	CMD_PWCTRL2   = 0xc1 // Power Control 2
	CMD_PWCTRL3   = 0xc2 // Power Control 3 (for Normal Mode)
	CMD_PWCTRL4   = 0xc3 // Power Control 4 (for Idle Mode)
	CMD_PWCTRL5   = 0xc4 // Power Control 5 (for Partial Mode)
	CMD_VMCTRL    = 0xc5 // VCOM Control
	CMD_CABCCTRL1 = 0xc6 // CABC Control 1
	CMD_CABCCTRL2 = 0xc8 // CABC Control 2
	CMD_CABCCTRL3 = 0xc9 // CABC Control 3
	CMD_CABCCTRL4 = 0xca // CABC Control 4
	CMD_CABCCTRL5 = 0xcb // CABC Control 5
	CMD_CABCCTRL6 = 0xcc // CABC Control 6
	CMD_CABCCTRL7 = 0xcd // CABC Control 7
	CMD_CABCCTRL8 = 0xce // CABC Control 8
	CMD_CABCCTRL9 = 0xcf // CABC Control 9

	CMD_NVMWR     = 0xd0 // NV Memory Write
	CMD_NVMPKEY   = 0xd1 // NV Memory Protection Key
	CMD_NVMSRD    = 0xd2 // NV Memory Status Read
	CMD_RDID4     = 0xd3 // Read ID4
	CMD_ADJCTRL1  = 0xd7 // Adjust Control 1
	CMD_PGAMCTRL  = 0xe0 // Positive Gamma Control
	CMD_NGAMCTRL  = 0xe1 // Negative Gamma Control
	CMD_DGAMCTRL1 = 0xe2 // Ditigal Gamma Control 1
	CMD_DGAMCTRL2 = 0xe3 // Ditigal Gamma Control 2
	CMD_SETIMAGE  = 0xe9 // Set Image Function
	CMD_ADJCTRL2  = 0xf2 // Adjust Control 2
	CMD_ADJCTRL3  = 0xf7 // Adjust Control 3
	CMD_ADJCTRL4  = 0xf8 // Adjust Control 4
	CMD_ADJCTRL5  = 0xf9 // Adjust Control 5
	CMD_SPIRDCMDS = 0xfb // SPI Read Command Setting
	CMD_ADJCTRL6  = 0xfc // Adjust Control 6
)

const (
	MADCTRL_MY  uint8 = 0x80 // Row Address Order         1 = address bottom to top
	MADCTRL_MX        = 0x40 // Column Address Order      1 = address right to left
	MADCTRL_MV        = 0x20 // Row/Column Exchange       1 = mirror and rotate 90 ccw
	MADCTRL_ML        = 0x10 // Vertical Refresh Order    1 = refresh bottom to top
	MADCTRL_BGR       = 0x08 // RGB-BGR Order             1 = Blue-Green-Red pixel order
	MADCTRL_MH        = 0x04 // Horizontal Refresh Order  1 = refresh right to left
)
