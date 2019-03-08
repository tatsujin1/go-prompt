package prompt

// Option is the type to replace default parameters.
// prompt.New accepts any number of options (this is functional option pattern).
type Option func(prompt *Prompt) error

// OptionParser to set a custom ConsoleParser object. An argument should implement ConsoleParser interface.
func OptionParser(x ConsoleParser) Option {
	return func(p *Prompt) error {
		p.in = x
		return nil
	}
}

// OptionWriter to set a custom ConsoleWriter object. An argument should implement ConsoleWriter interface.
func OptionWriter(x ConsoleWriter) Option {
	return func(p *Prompt) error {
		registerConsoleWriter(x)
		p.renderer.out = x
		return nil
	}
}

// OptionTitle to set title displayed at the header bar of terminal.
func OptionTitle(x string) Option {
	return func(p *Prompt) error {
		p.renderer.title = x
		return nil
	}
}

// OptionPrefix to set prefix string.
func OptionPrefix(x string) Option {
	return func(p *Prompt) error {
		p.renderer.prefix = x
		return nil
	}
}

// OptionCompletionWordSeparator to set word separators. Enable only ' ' if empty.
func OptionCompletionWordSeparator(x string) Option {
	return func(p *Prompt) error {
		p.completion.wordSeparator = x
		return nil
	}
}

// OptionLivePrefix to change the prefix dynamically by callback function
func OptionLivePrefix(f func() (prefix string, useLivePrefix bool)) Option {
	return func(p *Prompt) error {
		p.renderer.livePrefixCallback = f
		return nil
	}
}

// OptionPrefixTextColor change a text color of prefix string
func OptionPrefixTextColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.prefixText = x
		return nil
	}
}

// OptionPrefixBackgroundColor to change a background color of prefix string
func OptionPrefixBackgroundColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.prefixBG = x
		return nil
	}
}

// OptionInputTextColor to change a color of text which is input by user
func OptionInputTextColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.inputText = x
		return nil
	}
}

// OptionInputBGColor to change a color of background which is input by user
func OptionInputBGColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.inputBG = x
		return nil
	}
}

// OptionPreviewChoiceTextColor to change a text color which is completed
func OptionPreviewChoiceTextColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.previewChoiceText = x
		return nil
	}
}

// OptionPreviewChoiceBGColor to change a background color which is completed
func OptionPreviewChoiceBGColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.previewChoiceBG = x
		return nil
	}
}

// OptionChoiceTextColor to change a text color in drop down suggestions.
func OptionChoiceTextColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.choiceText = x
		return nil
	}
}

// OptionChoiceBGColor change a background color in drop down suggestions.
func OptionChoiceBGColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.choiceBG = x
		return nil
	}
}

// OptionSelectedChoiceTextColor to change a text color for completed text which is selected inside suggestions drop down box.
func OptionSelectedChoiceTextColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.selectedChoiceText = x
		return nil
	}
}

// OptionSelectedChoiceBGColor to change a background color for completed text which is selected inside suggestions drop down box.
func OptionSelectedChoiceBGColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.selectedChoiceBG = x
		return nil
	}
}

// OptionDescriptionTextColor to change a background color of description text in drop down suggestions.
func OptionDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.descriptionText = x
		return nil
	}
}

// OptionDescriptionBGColor to change a background color of description text in drop down suggestions.
func OptionDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.descriptionBG = x
		return nil
	}
}

// OptionSelectedDescriptionTextColor to change a text color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.selectedDescriptionText = x
		return nil
	}
}

// OptionSelectedDescriptionBGColor to change a background color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.selectedDescriptionBG = x
		return nil
	}
}

// OptionScrollbarThumbColor to change a thumb color on scrollbar.
func OptionScrollbarThumbColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.scrollbarThumb = x
		return nil
	}
}

// OptionScrollbarBGColor to change a background color of scrollbar.
func OptionScrollbarBGColor(x Color) Option {
	return func(p *Prompt) error {
		x, _ = p.renderer.ValidateColor(x)
		p.renderer.Colors.scrollbarBG = x
		return nil
	}
}

// OptionMaxVisibleChoices specify the max number of displayed completion choices.
func OptionMaxVisibleChoices(x uint16) Option {
	return func(p *Prompt) error {
		p.completion.maxVisibleChoices = int(x)
		return nil
	}
}

// OptionHistory to set history expressed by string array.
func OptionHistory(x []string) Option {
	return func(p *Prompt) error {
		p.history.histories = x
		p.history.Clear()
		return nil
	}
}

// OptionEditMode set a key bind mode.
func OptionEditMode(m EditMode) Option {
	return func(p *Prompt) error {
		p.editMode = m
		return nil
	}
}

// OptionBindKey to bind keys to functions
func OptionBindKey(b ...KeyBind) Option {
	return func(p *Prompt) error {
		for _, bind := range b {
			p.keyBindings[bind.Key] = bind.Fn
		}
		return nil
	}
}

// OptionBindControlSequence to make a binding to a specific control sequence.
func OptionBindControlSequence(b ...ControlSequenceBind) Option {
	return func(p *Prompt) error {
		for _, bind := range b {
			p.ControlSequenceBindings[bind.Sequence] = bind.Fn
		}
		return nil
	}
}

// OptionShowCompletionAtStart to set completion window is open at start.
func OptionShowCompletionAtStart() Option {
	return func(p *Prompt) error {
		p.completion.showAtStart = true
		return nil
	}
}

// New returns a Prompt with powerful auto-completion.
func New(executor Executor, completer Completer, opts ...Option) *Prompt {
	defaultWriter := NewStdoutWriter()
	registerConsoleWriter(defaultWriter)

	renderer := NewRender("> ", defaultWriter)

	pt := &Prompt{
		in:          NewStandardInputParser(),
		renderer:    renderer,
		buf:         NewBuffer(),
		executor:    executor,
		history:     NewHistory(),
		completion:  NewCompletionManager(completer, 6),
		editMode:    EmacsMode, // All the above assume that bash is running in the default Emacs setting
		keyBindings: make(map[KeyCode]KeyBindFunc, 10),
	}

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt
}
