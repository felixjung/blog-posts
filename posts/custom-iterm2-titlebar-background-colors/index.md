# Custom iTerm2 Titlebar Background Colors

If you are like me, you get pretty obsessed about the look and feel of your
development environment. It’s such a productivity waste, but I am a sucker for
nice looking apps. After all, as a dev, you look at this stuff all day.

Over the past years, Electron apps have become all the rage. If you look at a
typical JavaScript dev’s development environment these days, you will probably
notice them running either [Atom](https://atom.io) or
[Visual Studio Code](https://code.visualstudio.com) (VSCode) for their editor
and [Hyper](https://hyper.is) for their terminal emulator. In particular VS Code
seems to have become the default editor for most. Personally, I do not really
like the look of VS Code. Yes, it has made a lot of progress and there are some
cool themes, but to me, it still really looks like a Microsoft product. At least
as soon as you start digging a level deeper. Atom and Hyper on the other hand,
not so much. Look at these two beautiful apps next to each other. They are
🔥🔥🔥🔥.

![Atom on the left, Hyper on the right. Two beautiful Electron apps.](images/atom-and-iterm.jpg 'Screenshot showing the Atom text editor and an iTerm2 window next to each other.')

Sadly, though, I cannot get myself to use these apps for work. I am used to
[Neovim](https://neovim.io)/Vim and Atom’s model for splits is way too
different. It just does not feel right. And while Hyper is beautiful and its
plugin system rocks, as soon as I try to run Vim in it, I start noticing all
kinds of performance related glitches and lag. On top of all of that, these
Electron apps just take too much of my system’s resources. Today, I noticed Atom
causing 150% CPU load when I was scrolling its settings screen. _Scrolling_.
Both apps regularly take up gigabytes of memory.

So, I am stuck with my iTerm2, tmux, Neovim combo, which is not easy to get to
look really nice. With enough patience and determination, I think I ended up
with a nicely looking setup (getting theme colors to work has become a lot
easier thanks to [True Color](https://gist.github.com/XVilka/8346728) support
across the stack). You be the judge.

![Neovim on the left, terminal on the right, all running inside a tmux session in iTerm2.](DraggedImage.jpg)

However, one thing that always bothered me was iTerm’s grey titlebar. In the
nightly builds, George Nachman has added support for macOS’ dark mode for the
titlebar. This is nice, but still does not look as beautiful as Hyper.

![iTerm2 with dark titlebar left, Hyper on the right.](DraggedImage-1.jpg)

Recently, I finally — by accident — found a solution to this problem. 🎉 Turns
out, iTerm2 supports some
[proprietary escape codes](https://www.iterm2.com/documentation-escape-codes.html)
for overwriting certain settings. And yes, there are codes to change the
titlebar background color!

Here is how to set the titlebar color to the background of the
[iTerm2 One Dark](https://github.com/joshdick/onedark.vim/blob/master/term/One%20Dark.itermcolors)
theme I am using. Type the following into your iTerm2 terminal.

```bash
$ echo -e "\033]6;1;bg;red;brightness;40\a"
$ echo -e "\033]6;1;bg;green;brightness;44\a"
$ echo -e "\033]6;1;bg;blue;brightness;52\a"
```

This is what it looks like.

![iTerm2 with custom titlebar background.](DraggedImage-2.jpg)

Finally. iTerm will not save this to your settings. The solution is to
automatically execute the above commands when launching a new session. Just add
them to your .bashrc or .zshrc, etc.
