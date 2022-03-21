# macOS Linux mode

All-in-one project to help you get PC-like experience on your maOS. Mostly for developers & other IT guys, but not only
😉.

What you can get:
- PC keyboard shortcuts on your macOS, browser (chromium based) and IntelliJ IDEA
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


<table>
<tr>
<th align="center">
<img width="1000" height="1">
<p>
<small>
PC keyboard shortcut
</small>
</p>
</th>
<th align="center">
<img width="1000" height="1">
<p>
<small>
Mac keyboard shortcut
</small>
</p>
</th>
<th align="center">
<img height="1">
<p>
<small>
Notes
</small>
</p>
</th>
</tr>


<tr>
<td>
<kbd>Win</kbd>
</td>
<td>
<kbd>⌥</kbd>
</td>
<td>
Open Spotlight (can be easily changed to run _Launcher_ or similar apps)
</td>
</tr>


<tr>
<td>
<kbd>Ctrl</kbd> + <kbd>Alt</kbd> + <kbd>T</kbd>
</td>
<td>
<kbd>⌃</kbd> + <kbd>⌘</kbd> + <kbd>T</kbd>
</td>
<td>
Open iTerm
</td>
</tr>


<tr>
<td>
<kbd>Ctrl</kbd> + <kbd>Alt</kbd> + <kbd>S</kbd>
</td>
<td>
<kbd>⌃</kbd> + <kbd>⌘</kbd> + <kbd>S</kbd>
</td>
<td>
Open _Preferences_ in IntelliJ IDEA (workaround for https://youtrack.jetbrains.com/issue/IDEA-164155)
</td>
</tr>


<tr>
<td>
<kbd>Opt</kbd> & <kbd>Cmd</kbd> swap
</td>
<td>
-
</td>
<td>
Thanks to this rule the setup also works with Mac keyboards, thus it should be active only for Mac devices. To add or change supported devices check identifiers section in the rule and set vendor_id & product_id. Use Karabiner-EventViewer.app to check your device details.
</td>
</tr>


<tr>
<td>
<kbd>Alt</kbd> + ` (+ <kbd>4</kbd> / + <kbd>5</kbd> / +  <kbd>6</kbd>)
</td>
<td>
<kbd>⌘</kbd> + ` (+ <kbd>4</kbd> / + <kbd>5</kbd> / + <kbd>6</kbd>)
</td>
<td>
Workaround rules for VSC popup issue in IntelliJ IDEA (Show History, Git Blame and Show Diff)
</td>
</tr>
</table>

>Karabiner-Elements stores rules in under `~/.config/karabiner/assets/complex_modifications` path. To tweak
> the rules update the json files there and enable the rules in _Complex modifications_ tab.

>Bear in mind that many of the configured shortcuts may clash with the system ones, so you may need to disable some of them in the _System Preferences_.

>Rules order is important, remember about it about tweaking the existing ones and adding your own. 
### Importing IntelliJ IDEA keymap

1. Install IntelliJ plugin [XWin Keymap](https://plugins.jetbrains.com/plugin/13094-xwin-keymap) (it used to be preinstalled).
2. Copy [XWin IntelliJ IDEA.xml](https://github.com/raxigan/macos-linux-mode/blob/init/XWin%20IntelliJ%20IDEA.xml) file into the keymap configuration directory: `~/Library/Application Support/JetBrains/IntelliJIdea2021.3/keymaps` (the path may differ).
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

### What still works differently

TBD

## Dock and Switcher replacement

It is not my intention to hate the Dock here, but... Let's get rid of it 😄.

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

3. Install [AltTab](https://alt-tab-macos.netlify.app/). Example configuration:

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
- [@Damien](https://www.maketecheasier.com/author/damienoh/) for [Home & End keys fix](https://www.maketecheasier.com/fix-home-end-button-for-external-keyboard-mac/)
- [@Christian Long](https://apple.stackexchange.com/users/41838/christian-long) for [Dock auto-hide config](https://apple.stackexchange.com/a/82084)

## Contributing
TBD