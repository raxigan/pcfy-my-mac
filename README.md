# macOS Linux mode

All-in-one project to help you get PC-like experience on your maOS. Mostly for developers & other IT guys, but not only
😉.

What you can get:
- PC keyboard shortcuts on your macOS, browser (chromium based) and IntelliJ
    - the configuration works for both PC and Mac keyboards in same time (there is a special device-checking rule)
    - there are no custom shortcuts - all of them would do exactly the same on Linux/Windows or almost (e.g. <kbd>
      Win</kbd>/<kbd>Option</kbd> key opens Spotlight while on Windows or Linux Mint would open Start Menu)
- Dock & built-in switcher replacement
- Basic window management

## Keyboard shortcuts

### Importing karabiner-elements rules

1. Install [Karabiner-Elements](https://karabiner-elements.pqrs.org/)
2. Open Karabiner-Elements, create new profile (e.g. _Linux mode_) and select it:
<img src="./resources/karabiner-new-profile.png"/>

3. Open the following URL in your browser and allow the website to open Karabiner-Elements.app:

     ```
     karabiner://karabiner/assets/complex_modifications/import?url=https://raw.githubusercontent.com/raxigan/macos-linux-mode/init/linux-mode.json
     ```

4. Click _Import_:
<img src="./resources/karabiner-import.png"/>

5. Click _Enable All_:
<img src="./resources/karabiner-enable-all.png"/>

### Rules description (the less obvious ones)

| Rule   <img width=200/>                                                                        | Notes                                                                                                                                                                                                                                                                                       |
|--------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| <span><kbd>Win</kbd>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>                                                                 | Open Spotlight (can be easily changed to run _Launcher_ or similar apps)                                                                                                                                                                                                                    |
| <span><kbd>Ctrl</kbd> + <kbd>Alt</kbd> + <kbd>T</kbd></span>                                | Open iTerm                                                                                                                                                                                                                                                                                  |
| <kbd>Ctrl</kbd> + <kbd>Alt</kbd> + <kbd>S</kbd>                                | Open _Preferences_ in IntelliJ IDEA (workaround for https://youtrack.jetbrains.com/issue/IDEA-164155)                                                                                                                                                                                       |
| <kbd>Opt</kbd> & <kbd>Cmd</kbd> swap                                           | <p>Thanks to this rule the setup also works with Mac keyboards, thus it should be active only for Mac devices. To add or change supported devices check identifiers section in the rule and set _vendor_id_ & _product_id_. Use Karabiner-EventViewer.app to check your device details.</p> |
| <kbd>Alt</kbd> + ` (+ <kbd>4</kbd> / + <kbd>5</kbd> / +  <kbd>6</kbd>) | Workaround rules for VSC popup issue in IntelliJ IDEA (Show History, Git Blame and Show Diff)                                                                                                                                                                                               |

>Karabiner-Elements stores rules in under `~/.config/karabiner/assets/complex_modifications` path. To tweak
> the rules update the json files there and enable the rules in _Complex modifications_ tab.

>Bear in mind that many of the configured shortcuts may clash with the system ones, so you may need to disable some of them in the system.

>Rules order is important, remember about it about tweaking the existing ones and adding your own. 
### Importing IntelliJ IDEA keymap

Importing Karabiner-Elements beforehand is required.
1. Install IntelliJ plugin [XWin Keymap](https://plugins.jetbrains.com/plugin/13094-xwin-keymap) (it used to be preinstalled).
2. Copy [XWin IntelliJ.xml](https://github.com/raxigan/macos-linux-mode/blob/init/XWin%20IntelliJ.xml) file into the keymap configuration directory: `~/Library/Application Support/JetBrains/IntelliJIdea2021.3/keymaps` (the path may differ).
3. Restart IntelliJ IDEA and go to Preferences → Keymap and in the dropdown select *XWin IntelliJ* keymap.

> Some configured shortcuts cannot be changed in the IDE because of validation there. For such cases
> it's required to perform the tweaks directly in the keymap then restart IDE.

### Fixing <kbd>Home</kbd> & <kbd>End</kbd> keys

By default, macOS does not bind <kbd>Home</kbd> & <kbd>End</kbd> keys to any function. 
To make them to, respectively, move caret to beginning and end of line run the following command within
this project root directory (clone it first): 
```
mkdir ~/Library/KeyBindings && cp DefaultKeyBinding.dict ~/Library/KeyBindings
```

Then restart your Mac.

## Dock and Switcher replacement

There is no option to hide the Dock completely, so it's required to tweak its auto-hide configuration.

1. In macOS _System Preferences_ → _Dock & Menu Bar_ enable _Automatically hide and show the Dock_.
2. Configure the Dock to show after 2 second (in case you really need) it by running this command in your terminal:
```
defaults write com.apple.dock autohide-delay -float 2; killall Dock
```

> To restare default Dock settings run the command:
> ```
> defaults delete com.apple.dock autohide-delay; killall Dock
>```

3. Install [AltTab](https://alt-tab-macos.netlify.app/)

Example configuration:

<img src="./resources/alttab_controls.png"/>
<img src="./resources/alttab_appearance.png"/>

> To keep the AltTab list clean and short it's recommended to use
> the Blacklists feature to exclude less frequently used apps from it and/or
> the ones that can be easily accessed in some other way (e.g. Spotify, Docker Desktop etc.)

## Window management

Install [Rectangle](https://rectangleapp.com/), then if you want to operate using <kbd>Win</kbd> and arrows you may want to set it up the following way: 

<img src="./resources/rectangle_settings.png"/>

> Before setting up Rectangle shortcuts select _Default_ profile in Karabiner-Elements, then set up the shortcuts
> and back to your custom profile again.

## Credits
- [@rux616](https://github.com/rux616) for [karabiner-windows-mode](https://github.com/rux616/karabiner-windows-mode)
- [@tezeko](https://github.com/tekezo) for [Karabiner-Elements](https://github.com/pqrs-org/Karabiner-Elements)
- [@serhii-londar](https://github.com/serhii-londar) for [open-source-mac-os-apps](https://github.com/serhii-londar/open-source-mac-os-apps)
- [@\[)amien](https://damieng.com/blog/2015/04/24/make-home-end-keys-behave-like-windows-on-mac-os-x/) for Home & End keys fix
- [@Christian Long](https://apple.stackexchange.com/users/41838/christian-long) for Dock auto-hide config
## Contributing
TBD