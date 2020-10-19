# Upgrading to Webpack 2 — What I Learned Along the Way

I upgraded our Angular 1.6 app at work from webpack 1 to webpack 2, mainly for
two reasons:

1. We moved from a gulp based monolithic build to webpack in summer, 2015, and
   did not update the configuration much after. Since then, a lot of things have
   changed and new best practices have emerged. With the upgrade to webpack 2 I
   also wanted to optimize our build in ways that were already possible with
   webpack 1.
2. Webpack 2 enables you to apply a technique called tree-shaking, which should
   reduce the overall build size. Who would not want that?

This post documents the changes I made to our build and what I learned along the
way. I assume that you are familiar with webpack and some of its key concepts
like entries, loaders, plugins, and chunks. If not, I highly recommend
[SurviveJS — webpack](https://leanpub.com/survivejs-webpack) by
[Juho Vepsäläinen](https://github.com/bebraw), a webpack core team member. You
can also find a great getting started guide in the
[guides section](https://webpack.js.org/guides/get-started/) of the
[new webpack website](https://webpack.js.org). While I was working on this post,
[Adam Rackis](https://twitter.com/AdamRackis) wrote a
[great post](https://medium.com/@adamrackis/vendor-and-code-splitting-in-webpack-2-6376358f1923#.pyuadzo48)
on code splitting with webpack 2. It helped me improve our build even further.

## Why Upgrade to Webpack 2

Tobias Koppers (a.k.a [@sokra](https://github.com/sokra)), the creator of
Webpack, published a
[gist of changes](https://gist.github.com/sokra/27b24881210b56bbaff7) in webpack
2, some time ago. While this gist is not up-to-date, it should contain most
reasons why upgrading to webpack 2 might be a good idea.
[Sean T. Larkin](https://twitter.com/TheLarkInn) of the webpack core-team also
published
[this blog post](https://medium.com/webpack/webpack-2-and-beyond-40520af9067f#.y5489bb6z),
listing reasons to upgrade.

Webpack 2 does offer some nice features like support for the
[ECMAScript (ES) dynamic import spec](http://www.2ality.com/2017/01/import-operator.html)
using the `import()` operator. This can be seen as a less powerful, but
potentially more elegant, replacement for
[webpack’s `require.ensure()`](https://webpack.js.org/guides/code-splitting-require/).

The big new feature in Webpack 2, however, is native support for ES6 modules,
facilitating the above-mentioned _tree-shaking_ technique. I will discuss this
in more detail in a subsequent section. But, essentially, this allows webpack to
drop unused exports from ES6 modules. The result are smaller chunks and
**shorter load times**. 🚀

### Starting Point

Here is where I started from:

- Roughly 37,000 lines of JavaScript. Specifically, Angular 1.6 with a lot of
  inline templates, written in ES6.
- Heavy use of [lodash](https://lodash.com) throughout the project.
  Unfortunately, because our application grew over time, we are importing the
  main lodash object directly _everywhere_ (i.e., `import _ from 'lodash';`).
- A relatively basic Webpack 1.12 setup. We use a vendor chunk comprised of the
  various aliases specified in the Webpack config and split code from within the
  [Angular UI-router](https://github.com/angular-ui/ui-router/tree/master) using
  `require.ensure()`.

Here is what this looked like when analyzed using the _awesome_
[webpack-bundle-analyzer](https://github.com/yuffiy/webpack-bundle-analyzer)
plugin:

![Analysis of initial webpack chunks in the app.](images/bundle-initial.jpg 'Screenshot showing the initial webpack bundle analyzer output')

Angular 1.6 made up a huge part of our vendor chunk. This is something I would
not be able to change. The second biggest dependency in vendor was lodash. I
thought there should be a way to make this part smaller 🤔. Lastly, I had made
the mistake of importing all of [D3 version 4](https://github.com/d3/d3) in a
route of our app, despite it providing named member exports. Possibly something
I could change.

## Updating Our Config to the Webpack 2 Format

My first step was to update our configuration for webpack 2. The webpack team
has put together a [pretty good guide](https://webpack.js.org/guides/migrating/)
for this. It was very straight forward. A couple of notes on the experience.

- The new function format of the configuration is extremely useful. Being able
  to directly pass and evaluate environment variables is convenient.
- Personally, I found the new loaders configuration format much easier to
  understand.

Compare these two config extracts:

```js
// Old
module: {
  loaders: [
    {
      test: /\.css$/,
      loader: 'style-loader!css-loader?modules'
    }
  ];
}

// New
module: {
  rules: [
    {
      test: /\.css$/,
      use: [
        {
          loader: 'style-loader'
        },
        {
          loader: 'css-loader',
          options: {
            modules: true
          }
        }
      ]
    }
  ];
}
```

rules and use make much more sense to me than the old loaders and loader. Yes,
the old format also supported a loaders array — instead of the string example
above — where we now use use. But then you had loaders nested within loaders. 😐

## Cleaning Up Our Chunks

Alright, I was now running webpack 2. Next, I wanted to optimize the way the
application was split into chunks. In our config we were already using one
commons chunk for heavily used third-party dependencies. Specifically, these
were put into an
[explicit vendor chunk](https://webpack.js.org/plugins/commons-chunk-plugin/#explicit-vendor-chunk)
using the
[CommonsChunkPlugin](https://webpack.js.org/plugins/commons-chunk-plugin/). This
is probably what most people, in particular in situations without lazy loading
of additional chunks, end up with. At least, most tutorials I found end up
showing this solution. It can get tedious, as each new dependency has to be
added to the vendor chunk manually. In addition, the explicit vendor chunk does
not lend itself well to tree-shaking, because Webpack will end up placing each
vendor entry in the chunk in full. Luckily, the CommonsChunkPlugin gives you
very fine grained control over what to place in a commons chunk via the
minChunks option. minChunks can be a function that determines if a module should
be included in the commons chunk. It is passed two parameters; a module object
and a numeric value representing the number of chunks the module is included in.
I ended up configuring three commons chunks; one for application-wide
dependencies loaded at the application entry, one catch-all chunk including
dependencies required by child chunks of the entry chunk, and, finally, one
chunk for shared application modules loaded by child chunks.

### Vendor Chunk

Third-party dependencies usually have one thing in common: they are located
within `node_modules`. The module object passed to our minChunks function has a
property context. It is the path to the module file, excluding the file itself.
Hence, you can check if context contains the substring `node_modules` to include
a module in the common vendor chunk. This method is known as the
“[implicit common vendor chunk](https://webpack.js.org/guides/code-splitting-libraries/#implicit-common-vendor-chunk)”.
I used the following very standard configuration for the main vendor chunk:

```js
...,
plugins: [
  new webpack.optimize.CommonsChunkPlugin({
    name: 'vendor',
    minChunks: ({ context }) => context &&
      context.indexOf('node_modules') > -1
  })
],
...
```

### Catch-all Vendor Chunk

The problem with the above vendor chunk was that it would only capture
dependencies required in the application’s entry chunk. All child chunks would
still contain third-party dependencies they needed, unless they were already
included in the main vendor chunk. Webpack behaves like this because users
should not have to load code they might never need synchronously — and the main
vendor chunk has to be loaded synchronously. By adding an extra _catch-all_
commons chunk, you reduce duplicate code in the children. Here is how I added
such a catch-all chunk to our application.

```js
...
new webpack.optimize.CommonsChunkPlugin({
  filename: 'shared-vendor-deps.[chunkhash].js',
  async: 'shared-node-deps',
  minChunks: ({ context }, count) => context
    && context.indexOf('node_modules') > -1 && count >= 2
}),
...
```

By adding the async option to the chunk configuration, I was able to change both
the name of the chunk and cause it to be loaded asynchronously. Webpack ended up
loading the chunk with the application’s entry chunk, which is slightly
contradictory to the
[plugin’s documentation](https://webpack.js.org/plugins/commons-chunk-plugin/#options).
I will follow up on the various configurations for the CommonsChunkPlugin in
future posts.

### Shared App Modules Chunk

Our application has a lot of code that gets used in several places. For example,
our AddressService is used in our shop module to look up a shipping address.
Similarly, the service is used in the signup and account modules to provide
address auto-completion via Google maps. All three modules are, however, located
in different chunks, which are created by dynamic imports in the router logic.
As with the third-party dependencies required by multiple child chunks,
webpack’s default behavior in such a case is to include the AddressService in
each of the three child chunks, if it is not already included in the app’s entry
chunk. To create a commons chunk for shared application logic, I added another
instance of the CommonsChunkPlugin with the following configuration.

```js
...
new webpack.optimize.CommonsChunkPlugin({
  filename: 'shared-app-modules.[chunkhash].js',
  async: 'shared-app-modules',
  minChunks: ({ context }, count) => context
    && context.indexOf('app/scripts') > -1 && count >= 2
}),
...
```

The difference to the catch-all chunk above lies merely in the string
app/scripts, looked up in the module context. This is the relative path where
our application code lives.

### Chunky Results

Having reconfigured our chunks, I got the following output analysis from
webpack-bundle-analyzer.

![Analysis of our chunks after replacing the explicit with an implicit commons vendor chunk and introducing a chunk for shared app modules.](images/webpack-bundle-chunks-1.jpg 'Screenshot showing webpack bundle analyzer output')

As you can see, the new chunk configuration put most of our third-party
dependencies into the entry-vendor-deps commons chunk, effectively replacing the
old vendor commons chunk. Additionally, we see the newly created
shared-app-modules chunk with our shared application code.

You will notice that D3 and RxJS are still included in our split chunks, despite
being imported from `node_modules`. This is because they are only used in their
respective modules and not in the application entry. Additionally, at this
stage, the entry-vendor-deps chunk only contained code extracted from the entry
chunk.

## Preparing the Build for Tree-shaking

With our chunks in order, it was time to finally get our project into shape for
tree-shaking and dead code elimination. Tree-shaking is an optimization
technique that — essentially — removes any unused ES6 module exports from your
build output. In webpack 2, this is a two-step process. In the first step
webpack identifies any unused exports and removes their export statement (i.e.,
the actual _tree-shaking_). In the second step, these now dead entities are
removed by dead code elimination during minification with UglifyJS. You will
find more details on this in
[Axel Rauschmayer’s post](http://www.2ality.com/2015/12/webpack-tree-shaking.html).
According to Tobias Koppers’ aforementioned Gist, webpack is only able to detect
unused exports in a module in the following situations:

- The module’s members are imported via a named import.
- The module’s default export is imported.
- Some of the module’s members are re-exported

For me this translated to: “If you want to reduce chunk sizes via tree-shaking,
better use named member imports.“ Now, I couldn’t do anything about this with
the Angular 1.6 part of our codebase. But, I could do something about some other
big dependencies visible in the bundle analysis above.

Here is how we typically imported the heavily used third-party libraries in our
across our codebase before I made changes:

```js
/**
 * In our application entry where we bootstrap Angular.
 * Not suitable for tree-shaking.
 **/
import angular from 'angular';

/**
 * All across the application.
 * Not suitable for tree-shaking.
 **/
import _ from 'lodash';
import moment from 'moment';

/**
 * In a module requiring some reactive programming.
 * GREAT for tree-shaking.
 **/
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/merge';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/last';
import 'rxjs/add/operator/scan';

/**
 * In a module that provides some charts.
 * Not suitable for tree-shaking.
 **/
import d3 from 'd3';
```

### Importing Module Members by Name Where Possible

The only dependency that was already being imported in a tree-shaking compatible
way was RxJS — 👍🏻 for that. With the huge dependency that D3 was, according to
the above chunk analysis, I was happy to see it had supported named member
imports of individual modules since version 4. So I matched the specific
functions we were using with their respective modules. Given our limited use of
D3 this was luckily not too much work.

```
/**
 * We replaced our good old
 *
 * import * as d3 from 'd3';
 *
 * with what's below.
 **/

// Imports for a bar chart we made.
import { timeFormat } from 'd3-time-format';
import { timeDays, timeMondays, timeSundays } from 'd3-time';
import { select, selectAll, classed, attr, event } from 'd3-selection'; // eslint-disable-line no-unused-vars
import { extent } from 'd3-array';
import { axisBottom, axisLeft } from 'd3-axis';
import { scaleTime, scaleLinear } from 'd3-scale';
import { format } from 'd3-format';
import { transition, delay, duration } from 'd3-transition'; // eslint-disable-line no-unused-vars

// Some imports for a pie chart.
import { scaleOrdinal } from 'd3-scale';
import { select } from 'd3-selection';
import { transition } from 'd3-transition'; // eslint-disable-line no-unused-vars
import { arc, pie } from 'd3-shape';
import { interpolate } from 'd3-interpolate';
```

As you can see below, after replacing our D3 imports, I was left with a much
smaller footprint of the library in our build.

![Analysis of chunks after importing D3 by named members (red rectangle in the middle).](images/webpack-bundle-chunks-d3.jpg 'Screenshot of webpack bundle analyzer output')

The chunk containing D3 (the big purple one) went from second- to third-largest.
Sweet 🍭.

### Handling a Monolithic Lodash Dependency

Now to the — hopefully — final step before flipping the switch enabling
tree-shaking. Finding a way to import only the bits of lodash we needed. Lodash
is available in many different builds. There is regular
[lodash](https://www.npmjs.com/package/lodash) (what we use),
[lodash-es](https://www.npmjs.com/package/lodash-es) (an ES2015 exports based
version),
[lodash-modularized](https://www.npmjs.com/browse/keyword/lodash-modularized)
(individual lodash functions as modules). You can also directly import
individual lodash functions (e.g., import map from 'lodash/map';). So, there are
ways to facilitate named member imports from lodash, enabling tree-shaking along
the way. However, I did not want to go through our entire codebase and rewrite
all those lodash imports. Luckily, there is a very handy
[lodash Babel plugin](https://github.com/lodash/babel-plugin-lodash) that turns
any imports of lodash and subsequent use of lodash methods into named member
imports. Here is the example from the plugin’s readme:

```js
// Turns
import _ from 'lodash';
import { add } from 'lodash/fp';

const addOne = add(1);
_.map([1, 2, 3], addOne);

// Roughly to
import _add from 'lodash/fp/add';
import _map from 'lodash/map';

const addOne = _add(1);
_map([1, 2, 3], addOne);
```

Without touching any of our code, I could get function imports for the entire
project. Using the plugin is very much straight forward. After installation, all
I did was add it to the list of babel plugins.

```js
use: [
  {
    loader: 'babel-loader',
    options: {
      plugins: ['lodash', 'syntax-dynamic-import'],
      presets: ['latest'],
      cacheDirectory: true,
      babelrc: false
    }
  }
];
```

And here is what happened with our chunks afterwards.

![All lodash dependencies are included in the entry-vendor-deps and shared-node-deps (bottom-left) chunks.](images/webpack-bundle-lodash.jpg 'Screenshot of webpack bundle analyzer output')

As illustrated by the many tiny boxes inside the lodash box in entry-vendor-deps
chunk above, enabling the lodash Babel plugin _did_ work as expected. You can
also see the _catch-all_ chunk caught all the child chunks’ lodash imports
(bottom-right). Awesome 🦄. There is a way to further optimize lodash’s
footprint with the
[lodash-webpack-plugin](https://github.com/lodash/lodash-webpack-plugin). As I
understand it, this plugin allows you to disable parts of lodash entirely,
causing them to be dropped from the build.

## Enabling Tree-shaking

And now for the grand finale. I had done about as much as I could to prepare our
codebase for tree-shaking with webpack 2. If not, please let me know! All that
was left was enabling the feature by

1. telling Babel to not transpile ES6 modules into CommonJS modules and
2. making sure JavaScript minification with the UglifyJS plugin was enabled.

The first part is necessary for webpack to be able to perform the static module
analysis and remove export statements from unused members. The second step
eliminates the resulting dead code.

To tell Babel not to transpile modules, set your ES preset’s modules option to
false:

```js
use: [
  {
    loader: 'babel-loader',
    options: {
      plugins: ['lodash', 'syntax-dynamic-import'],
      presets: [['latest', { modules: false }]],
      cacheDirectory: true,
      babelrc: false
    }
  }
];
```

Notice the nested array structure for the presets option. Because you may pass
multiple presets to Babel, you need to contain each preset with options in an
array of its own. There used to be a es2015-native-modules preset. This has been
deprecated in favor of the modules option.

There are (at least) two ways of enabling JavaScript minification in webpack.
First, you can pass the --optimize-minimize flag to the CLI. This will enable
the UglifyJsPlugin bundled with webpack. Second, you can add the UglifyJsPlugin
to the plugins in your configuration file. This also allows you to use a version
different from the one bundled with webpack. Here is the configuration I ended
up using for our production build:

```js
plugins: [
  ...,
  new webpack.optimize.UglifyJsPlugin({
    compress: {
      warnings: false,
      conditionals: true,
      unused: true,
      comparisons: true,
      sequences: true,
      dead_code: true,
      evaluate: true,
      if_return: true,
      join_vars: true
    },
    output: { comments: false }
  }),
  ...
]
```

Notice the `dead_code` option, which enables dead code elimination.

After all this and some more work in other parts, here is what our final
production build chunks now looked like.

![Chunks with tree-shaking enabled.](images/webpack-bundle-tree-shaking.jpg 'Screenshot of webpack bundle analyzer output')

Admittedly, the chunks do not look any different compared to the previous step.
I have not seen any effect of tree-shaking, yet. Some initial investigation
showed I had brought down the overall size of our build, gzipped, and including
all child chunks from 884Kb to 520Kb. That is quite a lot 💪🏻. At the same time,
however, due to the asynchronously loaded shared-app-modules chunk, the
initially loaded code of our login screen had increased from 173Kb to 251Kb. So,
clearly I still need to optimize some things and do a complete analysis of the
new build’s performance. Results will be published in the near future.

## (Preliminary) Conclusion

I was amazed by all the possibilities webpack offers these days. Once you’re
willing to deep-dive into the loaders, plugins and configuration options you
realize there is _a lot_ you can fine-tune about your build. I am sure this is
all old news to people who have been riding the React train for some long time.
To me this was still an exciting experience coming from our old build
configuration. Let’s see what the webpack team will produce over the coming
months. 🤓
