package cli_test

import (
	"fmt"
	"github.com/silvertern/cli"
	"github.com/silvertern/cli/command"
	"github.com/silvertern/cli/option"
	"sort"
	"testing"
)

func setup_parse_app() cli.App {
	co := command.New("checkout", "checkout a branch or revision").
		WithShortcut("co").
		WithArg(command.Arg{Key: "branch"}).
		WithOption(option.New("branch", "Create branch").WithChar('b').WithType(option.TypeBool)).
		WithOption(option.New("upstream", "Set upstream").WithChar('u').WithType(option.TypeBool)).
		WithOption(option.New("fallback", "Set upstream").WithChar('f')).
		WithOption(option.New("count", "Count").WithChar('c').WithType(option.TypeInt)).
		WithOption(option.New("pi", "Set upstream").WithChar('p').WithType(option.TypeNumber)).
		WithOption(option.New("str", "Count").WithChar('s'))

	add := command.New("add", "add a remote").
		WithArg(command.Arg{Key: "remote"}).
		WithArg(command.Arg{Key: "count", Type: option.TypeInt}).
		WithArg(command.Arg{Key: "pi", Type: option.TypeNumber}).
		WithArg(command.Arg{Key: "force", Type: option.TypeBool}).
		WithArg(command.Arg{Key: "optional", Optional: true, Type: option.TypeBool}).
		WithOption(option.New("force", "Force").WithChar('f').WithType(option.TypeBool)).
		WithOption(option.New("quiet", "Quiet").WithChar('q').WithType(option.TypeBool)).
		WithOption(option.New("default", "Default"))

	rmt := command.New("remote", "operations with remotes").WithCommand(add)

	return cli.New("git tool").
		WithArg(command.Arg{Key: "arg1"}).
		WithCommand(co).
		WithCommand(rmt)
}

func TestApp_Parse_DropsPathFromAppName_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"~/some/path/git", "checkout", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[]", invocation, args, opts, err)
}

func TestApp_Parse_DropsDotPathFromAppName_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"./git", "checkout", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[]", invocation, args, opts, err)
}

func TestApp_Parse_NoFlags_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[]", invocation, args, opts, err)
}

func TestApp_Parse_1xCharBoolFlag_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-b", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true]", invocation, args, opts, err)
}

func TestApp_Parse_2xCharBoolFlags_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-b", "-u", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_2xCharBoolFlagsAsOne_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-bu", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_MultiCharStringLast_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-buf", "master", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true fallback:master upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_MultiCharIntLast_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-buc", "1", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true count:1 upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_MultiCharNumberLast_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-bup", "3.14", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true pi:3.14 upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_1xBoolFlag_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "--branch", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true]", invocation, args, opts, err)
}

func TestApp_Parse_2xBoolFlag_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "--branch", "--upstream", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_2xBoolAnd1xStringFlag_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "--fallback=master", "--branch", "--upstream", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true fallback:master upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_RedundantFlags_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-b", "--branch", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true]", invocation, args, opts, err)
}

func TestApp_Parse_NestedCommandWithFlags_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "-f", "1", "3.14", "true"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true] map[force:true]", invocation, args, opts, err)
}

func TestApp_Parse_OptionalMissing_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "1", "3.14", "true"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true] map[]", invocation, args, opts, err)
}

func TestApp_Parse_OptionalPresent_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "1", "3.14", "true", "true"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true true] map[]", invocation, args, opts, err)
}

func TestApp_Parse_KeysAnywhereBetweenArgs_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "-f", "origin", "--default=foo", "1", "3.14", "true", "-q"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true] map[default:foo force:true quiet:true]", invocation, args, opts, err)
}

func TestApp_Parse_ExplicitValueForBoolOption_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "--force=true", "origin", "1", "3.14", "true"})
	assertAppParseError(t, "[git remote add] [] map[]",
		"boolean options have true assigned implicitly, found value for --force", invocation, args, opts, err)
}

func TestApp_Parse_EqSignInStringOptionValue_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "--default=foo=boo=blah", "origin", "1", "3.14", "true"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true] map[default:foo=boo=blah]", invocation, args, opts, err)
}

func TestApp_Parse_UnrecognizedCommand_ErrorUnknownFlagForRoot(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "foo", "-f", "origin"})
	assertAppParseError(t, "[git] [foo] map[]", "unknown flag -f", invocation, args, opts, err)
}

func TestApp_Parse_UnrecognizedCommand_ErrorUnknownArgument(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "foo", "origin"})
	assertAppParseError(t, "[git] [foo origin] map[]", "unknown arguments [origin]", invocation, args, opts, err)
}

func TestApp_Parse_BaseApp_ErrorMissingArgument(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git"})
	assertAppParseError(t, "[git] [] map[]", "missing required argument arg1", invocation, args, opts, err)
}

func TestApp_Parse_DanglingOptions_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "dev", "-p"})
	assertAppParseError(t, "[git checkout] [dev] map[]", "dangling option --pi", invocation, args, opts, err)
}

func TestApp_Parse_LastArgOptionalRequiredMissing_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "1"})
	assertAppParseError(t, "[git remote add] [origin 1] map[]", "missing required argument pi", invocation, args, opts, err)
}

func TestApp_Parse_IncorrectBoolArgType_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "1", "3.14", "foo"})
	assertAppParseError(t, "[git remote add] [origin 1 3.14 foo] map[]", "argument force must be a boolean value, found foo", invocation, args, opts, err)
}

func TestApp_Parse_IncorrectIntArgType_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "3.14", "3.14", "true"})
	assertAppParseError(t, "[git remote add] [origin 3.14 3.14 true] map[]", "argument count must be an integer value, found 3.14", invocation, args, opts, err)
}

func TestApp_Parse_IncorrectNumberArgType_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "1", "aaa", "true"})
	assertAppParseError(t, "[git remote add] [origin 1 aaa true] map[]", "argument pi must be a number, found aaa", invocation, args, opts, err)
}

func TestApp_Parse_IncorrectOptionalArgType_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "1", "3.14", "true", "25"})
	assertAppParseError(t, "[git remote add] [origin 1 3.14 true 25] map[]", "argument optional must be a boolean value, found 25", invocation, args, opts, err)
}

func TestApp_Parse_NonBooleanFlagInNonTerminalPosition_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-bpu", "3.14", "dev"})
	assertAppParseError(t, "[git checkout] [] map[branch:true]", "non-boolean flag -p in non-terminal position", invocation, args, opts, err)
}

func TestApp_Parse_MissingValueForOption_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "--pi", "dev"})
	assertAppParseError(t, "[git checkout] [] map[]", "missing value for option --pi", invocation, args, opts, err)
}

func TestApp_Parse_NoValueAfterTheEqualSignForStringOption_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "--str=", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[str:]", invocation, args, opts, err)
}

func TestApp_Parse_UnknownOption_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "--foo=25", "dev"})
	assertAppParseError(t, "[git checkout] [] map[]", "unknown option --foo", invocation, args, opts, err)
}

func TestApp_Parse_IncorrectDataForIntOption_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-c", "2.25", "dev"})
	assertAppParseError(t, "[git checkout] [dev] map[count:2.25]", "option --count must be given an integer value, found 2.25", invocation, args, opts, err)
}

func TestApp_Parse_IncorrectDataForNumberOption_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-p", "aaa", "dev"})
	assertAppParseError(t, "[git checkout] [dev] map[pi:aaa]", "option --pi must must be given a number, found aaa", invocation, args, opts, err)
}

func TestApp_Parse_LastArgOptionalPermitsUnlimitedExtraArgs_Error(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "remote", "add", "origin", "1", "3.14", "true", "true", "exra1", "extra2"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true true exra1 extra2] map[]", invocation, args, opts, err)
}

func TestApp_Parse_HelpOptionComesOutWithoutArgOrFlagValidation_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-c", "3.15", "dev", "arg2", "arg3", "--help"})
	assertAppParseOk(t, "[git checkout] [] map[help:true]", invocation, args, opts, err)
}

func TestApp_Parse_HelpFlagInMultichar_Ok(t *testing.T) {
	invocation, args, opts, err := cli.Parse(setup_parse_app(), []string{"git", "checkout", "-bhc", "3.15", "dev"})
	assertAppParseOk(t, "[git checkout] [] map[help:true]", invocation, args, opts, err)
}

func assertAppParseOk(t *testing.T, expected string, invocation []string, args []string, opts map[string]string, err error) {
	if err == nil {
		optkeys := []string{}
		for key := range opts {
			optkeys = append(optkeys, key)
		}
		sort.Strings(optkeys)
		for i, key := range optkeys {
			optkeys[i] = fmt.Sprintf("%s:%s", key, opts[key])
		}
		actual := fmt.Sprintf("%v %v map%v", invocation, args, optkeys)
		if actual != expected {
			t.Errorf("assertion error: expected '%v', found '%v'", expected, actual)
		}
	} else {
		t.Errorf("no error expected, found '%v'; data %v %v %v", err, invocation, args, opts)
	}
}

func assertAppParseError(t *testing.T, expectedData, expectedError string, invocation []string, args []string, opts map[string]string, err error) {
	optkeys := []string{}
	for key := range opts {
		optkeys = append(optkeys, key)
	}
	sort.Strings(optkeys)
	for i, key := range optkeys {
		optkeys[i] = fmt.Sprintf("%s:%s", key, opts[key])
	}
	actual := fmt.Sprintf("%v %v map%v", invocation, args, optkeys)
	if actual != expectedData {
		t.Errorf("assertion error: expectedData '%v', found '%v'", expectedData, actual)
	}
	if err == nil {
		t.Error("an error was expected")
	} else if expectedError != err.Error() {
		t.Errorf("error mismatch, expected: '%v', found '%v'", expectedError, err.Error())
	}
}