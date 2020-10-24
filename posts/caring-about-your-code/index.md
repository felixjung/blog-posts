# Caring About Your Code

When you are a developer and hacking together some scripts for yourself, there
is no immediate need to care about the quality of your code or its style.
Likewise when you are using Git to maintain some small personal project, why
should you are about things like commit messages or even pull request titles —
who creates branches in personal projects anyway, right?

It is easy to see it this way and keep on hacking together a bunch of code. Keep
on committing with messages like “Bug fixes”, ”More bug fixes”, ”More bug fixes
2”. _Every_. _Single_. _Time_. After all, it is just some side project, right? I
think most developers, myself included, have told themselves something along the
lines of “Ah come on, who cares. I’m just going to commit it like this.” before.
But you know what? I strongly believe every time we do that, a tiny part of the
better developer most of us want to become goes into a corner and dies. 💀

A big part of becoming a better developer is experience, sure. Write enough
lines of code in different languages and for enough problems and you will have a
pretty good idea of what you need to do most of the time. But this is only part
of the story. Another huge part, maybe equally as important, is caring about how
you write all that code; how you commit your code, too. Whenever you write code,
at every step, ask yourself: “Is there a better way to do this?”. Some examples:

- “Does this function or variable name communicate its meaning as precisely as
  possible?”
- “Instead of using this if-else block, can I maybe just use an if block and
  return early?”
- “Can I generalize this bit of code and make it a function instead of repeating
  myself all over the place?”
- “Is the piece of code I’m about to commit an atomic unit? Should I split it
  into multiple commits?”
- “Is my commit message easily understandable? Does it accurately summarize the
  changes in this commit? Can I make it even more concise?”

Always challenge yourself when you code. You might not find the perfect answer
to the above questions right away, but asking them and doing your best to
improve will eventually make you a better developer.

But it does not stop there. Not only will you become a better developer, other
contributors to that open source project you have been helping out on and your
colleagues at work will benefit, too! 🙋🙋🏾‍♂️ More examples:

- The junior your team just hired will be able to read your code and learn from
  it more easily!
- When you go on vacation and people find a bug in the code you committed just
  before leaving, they can quickly identify the commit that introduced it
  (remember, concise commit messages) and revert it (you made a nice atomic
  commit). You can take care of the problem once you get back.
- And that huge new feature branch you had been working on for months before it
  was squash-merged can also be reverted because you named your pull request
  nicely and even included the ticket number from whatever ticket system you
  poor soul have to use at your company — yes Jira, I am talking about you.

So yeah, just give a damn about this stuff. You will be a better developer for
it and people will enjoy working with you. Everybody wins. 🎉 Pro tip: this is
true for most things in life.
