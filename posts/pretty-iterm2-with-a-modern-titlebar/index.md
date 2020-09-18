# Pretty iTerm2 With a Modern Titlebar 💄💅

In a previous post I wrote about customizing iTerm2's titlebar to match your current theme. This provided a pretty neat visual upgrade and brought the user experience closer to that of the all hyped and very beautiful Hyper terminal emulator bei Zeit.co.

Unfortunately, [an API change in macOS High Sierra](https://gitlab.com/gnachman/iterm2/issues/4080#note_43758478) stopped the titlebar customization from working. [The author of iTerm2 had promised](https://gitlab.com/gnachman/iterm2/issues/4080#note_64566741) to work on changes bringing the feature back, but for a long time nothing seemed to happen. [That has changed now](https://gitlab.com/gnachman/iterm2/issues/4080#note_93855327).

A new combination of features and settings in iTerm2's nightly build allows you to configure the most beautiful terminal on macOS, yet. You will need to be on macOS Mojave to be able to do so.

![iTerm2 final setup](images/iterm2-final.tiff "Screenshot of iTerm2 window")

Here are the steps you need to take to get to this set up.

1. **Upgrade to macOS Mojave, if you have not already.**
   The upgrade went flawlessly for me and I have seen only positive reports online. To get started, open the App Store on your Mac, download the macOS Mojave installer, and follow the instructions.
2. **Install the iTerm2 nightly build.**
   I usually run the nightly on my systems. It is great and rarely causes problems.
3. **Set the iTerm2 theme to *Minimal*.**
   Go to *Settings*, select the *Appearance* tab, and choose *Minimal* from the *Theme* dropdown.
4. **Set your Profile's window style to *Compact*.**
   In the *Settings* go to the *Profiles* tab, choose *Window*, and in the *Settings for New Windows* section select *Compact* from the *Style* dropdown.
5. **Open a new terminal window.**

At this point, the new terminal window should look fairly similar to the above screenshot, matching your current theme. There are some additional *Advanced* settings of iTerm2, that you can tweak to get even closer or adjust things to your liking. Here's how I have set them.

![iTerm2 final setup advanced settings.](images/iterm2-advanced-settings.tiff "Screenshot of iTerm2 advanced settings")

That's it! :tada: Enjoy an amazingly beautiful terminal experience for endless developer Zen. Looking at the [issue on Gitlab](https://gitlab.com/gnachman/iterm2/issues/4080), this feature is not done yet and the experience might get even better. So stay tuned! Thanks so much to George Nachman for making such a great terminal emulator.

