# Test-driving Prepack with webpack 2 

About a month ago, Facebook announced [Prepack](https://prepack.io), a build tool that optimizes JavaScript for faster runtime execution. It does so by replacing any computations known at “compile time” with their result. Here’s an example from the Prepack webiste:

```js
// Input
(function () {
  function hello() { return 'hello'; }
  function world() { return 'world'; }
  global.s = hello() + ' ' + world();
})();

// Output
(function () {
  s = "hello world";
})();
```

Prepack still has a very long way to go before becoming production-ready for your average single-page application code. For example, at the moment, it is completely unaware of its environment. In the browser, document and window would evaluate to undefined, if you do not provide what the Prepack folks call a “model”.

However, that did not stop my colleagues at [@SumUpEng](https://twitter.com/sumupeng) and myself from taking Prepack for a spin in our Dashboard application — just to see what it would do in terms of performance.

### Setting up Prepack with webpack 2

Turns out, there is a [webpack 2 plugin](https://github.com/gajus/prepack-webpack-plugin) that came out shortly after Prepack was announced. All you have to do is to instantiate the plugin with a configuration object (possibly empty) and add it to the plugins property of your webpack config.

```js
import PrepackWebpackPlugin from ‘prepack-webpack-plugin’;

module.exports = (env) => ({
  // ...
  plugins: [
    new PrepackWebpackPlugin({})
  ]
});
```

One thing to note is that this plugin [will not work properly](https://github.com/gajus/prepack-webpack-plugin/issues/11) when you run it after the UglifyJS plugin. So make sure to position it before UglifyJS in the plugins array.

### What we learned

After doing some *very* casual tests loading a build of the Dashboard without and one with the Prepack webpack plugin, I was somewhat surprised. The overall size of the application had gone from 4.9MB for 5MB. Say what? I thought this was supposed to optimize our app? Why did it get bigger?

Luckily, my colleagues and I were at [@jsconfeu](https://twitter.com/jsconfeu?lang=en) when we made this experiment. So I walked over to Facebook’s [Christoph Pojer](https://twitter.com/cpojer) and some of his fellow engineers to figure out what was going on. And while these guys do not work on Prepack, they were quick to help me overcome some misconceptions I had.

Being a tool for optimizing runtime performance, Prepack does not care about file size! The only thing it should care about is making the code run faster inside the JavaScript VM. This means: simple instructions (i.e. a lot of simple assignments) and — in some cases — long computation results instead of the original computation code. Here is an example where Prepack will create significantly more code. You can go to the [interactive repl](https://prepack.io/repl.html) on prepack.io to try it out live.

```js
/**
 * Input @ 262 characters, measured by running
 * `wc -c prepack-example.js`.
 **/

const names = [
  'Han', 'Leia', 'Luke', 'Obi-Wan',
  'Anakin', 'Jyn', 'Rey', 'Finn'
]

function sayHello(name) {
  return `Hello ${name}! It is very nice to meet you. How are you doing?`
}

const greetAll = names.map(sayHello).join('\n')

console.log(greetAll)

/*
 * Prepack output @ 663 characters, measured by running
 *
 * prepack hello-star-wars.js \
 *   --out hello-star-wars-prepacked.js \
 *   && wc -c hello-star-wars-prepacked.js
 *
 **/

var sayHello;

(function () {
  function _0(name) {
    return `Hello ${name}! It is very nice to meet you. How are you doing?`;
  }

sayHello = _0;
  console.log("Hello Han! It is very nice to meet you. How are you doing?\nHello Leia! It is very nice to meet you. How are you doing?\nHello Luke! It is very nice to meet you. How are you doing?\nHello Obi-Wan! It is very nice to meet you. How are you doing?\nHello Anakin! It is very nice to meet you. How are you doing?\nHello Jyn! It is very nice to meet you. How are you doing?\nHello Rey! It is very nice to meet you. How are you doing?\nHello Finn! It is very nice to meet you. How are you doing?");
})();
```

Optimizing the example code using Prepack, we end up with a wooping **663** characters instead of the original 262.

> The time the VM spends parsing your bundle before anything happens seems like a big topic right now.

What does this mean for your application? Should you use Prepack or not, once it becomes ready for production? The answer will probably end up being: “it depends”. At least that’s what the conversation with [Christoph Pojer](https://twitter.com/cpojer) and colleagues came down to.

Here is why. The size of your bundle is most likely the biggest factor in how quickly your users will be able to interact with the page. However, the time the browser spends parsing your code before it can execute might add significantly to that time. Parsing performance seems like a big topic right now. [Addy Osmani](https://twitter.com/addyosmani)
[touched on it](https://youtu.be/7vUs5yOuv-o?t=23m10s) at JSConf EU and [Marja Hölttä](https://twitter.com/marjakh) [gave an entire talk](https://youtu.be/Fg7niTmNNLg) about how Chrome parses JavaScript. If the parsing time you save by optimizing with Prepack outweighs any additional load time due to an increased bundle size, it is worth using Prepack. Otherwise, probably not.

Personally, I was not aware of these performance implications. And having focussed so much on chunk sizes in the build output (see my post on [upgrading to webpack 2](https://codematters.blog/upgrading-to-webpack-2-fc09bd8adbd4) for more details on the topic), I had totally missed the point of Prepack. After some conversations on the topic with other developers, it seems I was not the only one, though.