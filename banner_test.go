package banner_test

import (
	"fmt"
	"testing"

	"github.com/j-edgizer/banner"
	"github.com/j-edgizer/banner/fonts"
)

func ExampleInline() {
	fmt.Println("start of banner")
	fmt.Println(banner.Inline("hey world.", fonts.Small))
	fmt.Println("end of banner")
	// Output:
	// start of banner
	//  _                                    _     _
	// | |_   ___  _  _   __ __ __ ___  _ _ | | __| |
	// | ' \ / -_)| || |  \ V  V // _ \| '_|| |/ _` | _
	// |_||_|\___| \_, |   \_/\_/ \___/|_|  |_|\__,_|(_)
	//             |__/
	// end of banner
}

func ExampleInline_lowercase() {
	fmt.Println("start of banner")
	fmt.Println(banner.Inline("abcdefghij", fonts.Small))
	fmt.Println(banner.Inline("klmnopqrst", fonts.Small))
	fmt.Println(banner.Inline("uvwxyz", fonts.Small))
	fmt.Println("end of banner")
	// Output:
	// start of banner
	//        _            _        __        _     _    _
	//  __ _ | |__  __  __| | ___  / _| __ _ | |_  (_)  (_)
	// / _` || '_ \/ _|/ _` |/ -_)|  _|/ _` || ' \ | |  | |
	// \__,_||_.__/\__|\__,_|\___||_|  \__, ||_||_||_| _/ |
	//                                 |___/          |__/
	//  _    _                                         _
	// | |__| | _ __   _ _   ___  _ __  __ _  _ _  ___| |_
	// | / /| || '  \ | ' \ / _ \| '_ \/ _` || '_|(_-<|  _|
	// |_\_\|_||_|_|_||_||_|\___/| .__/\__, ||_|  /__/ \__|
	//                           |_|      |_|
	//  _  _ __ ____ __ ____ __ _  _  ___
	// | || |\ V /\ V  V /\ \ /| || ||_ /
	//  \_,_| \_/  \_/\_/ /_\_\ \_, |/__|
	//                          |__/
	// end of banner
}

func ExampleInline_uppercase() {
	fmt.Println("start of banner")
	fmt.Println(banner.Inline("ABCDEFGHIJ", fonts.Small))
	fmt.Println(banner.Inline("KLMNOPQRST", fonts.Small))
	fmt.Println(banner.Inline("UVWXYZ", fonts.Small))
	fmt.Println("end of banner")
	// Output:
	// start of banner
	//    _    ___   ___  ___   ___  ___   ___  _  _  ___     _
	//   /_\  | _ ) / __||   \ | __|| __| / __|| || ||_ _| _ | |
	//  / _ \ | _ \| (__ | |) || _| | _| | (_ || __ | | | | || |
	// /_/ \_\|___/ \___||___/ |___||_|   \___||_||_||___| \__/
	//  _  __ _     __  __  _  _   ___   ___   ___   ___  ___  _____
	// | |/ /| |   |  \/  || \| | / _ \ | _ \ / _ \ | _ \/ __||_   _|
	// | ' < | |__ | |\/| || .` || (_) ||  _/| (_) ||   /\__ \  | |
	// |_|\_\|____||_|  |_||_|\_| \___/ |_|   \__\_\|_|_\|___/  |_|
	//  _   _ __   ____      ____  ____   __ ____
	// | | | |\ \ / /\ \    / /\ \/ /\ \ / /|_  /
	// | |_| | \ V /  \ \/\/ /  >  <  \ V /  / /
	//  \___/   \_/    \_/\_/  /_/\_\  |_|  /___|
	// end of banner
}

func TestInlineSmall(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"jjj", `
   _    _    _
  (_)  (_)  (_)
  | |  | |  | |
 _/ | _/ | _/ |
|__/ |__/ |__/
`},
		{"j j", `
   _      _
  (_)    (_)
  | |    | |
 _/ |   _/ |
|__/   |__/
`},
		{"j", `
   _
  (_)
  | |
 _/ |
|__/
`},
		{"@?!", `
  ____   ___  _
 / __ \ |__ \| |
/ / _` + "`" + ` |  /_/|_|
\ \__,_| (_) (_)
 \____/
`},
		{"ccc", `
 __  __  __
/ _|/ _|/ _|
\__|\__|\__|
`},
		{" ", `

`},
		{"", `

`},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := banner.Inline(test.input, fonts.Small)
			expected := test.expected[1 : len(test.expected)-1]
			if expected != output {
				t.Log("output: \n" + output)
				t.Log("expected: \n" + expected)
				t.Errorf("output differs")
			}
		})
	}
}

func TestInlineBanner(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"jjj", `
      #      #      #
      #      #      #
      #      #      #
      #      #      #
 #    # #    # #    #
  ####   ####   ####
`},
		{"j j", `
      #        #
      #        #
      #        #
      #        #
 #    #   #    #
  ####     ####
`},
		{"j", `
      #
      #
      #
      #
 #    #
  ####
`},
		{"@?!", `
  #####   #####  ###
 #     # #     # ###
 # ### #       # ###
 # ### #    ###   #
 # ####     #
 #               ###
  #####     #    ###
`},
		{"ccc", `
  ####   ####   ####
 #    # #    # #    #
 #      #      #
 #      #      #
 #    # #    # #    #
  ####   ####   ####
`},
		{" ", `

`},
		{"", `

`},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := banner.Inline(test.input, fonts.Banner)
			expected := test.expected[1 : len(test.expected)-1]
			if expected != output {
				t.Log("output: \n" + output)
				t.Log("expected: \n" + expected)
				t.Errorf("output differs")
			}
		})
	}
}
