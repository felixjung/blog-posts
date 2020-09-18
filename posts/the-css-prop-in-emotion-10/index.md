# The `css` Prop in Emotion 10

On November 27th, [Mitchell Hamilton](https://twitter.com/mitchellhamiltn) announced the release of Emotion 10 [in a blog post](https://medium.com/emotion-js/announcing-emotion-10-f1a4b17b8ccd). The new version comes with a lot of changes, a new API, and new packages. I encourage you to read the post and check out [the documentation](https://emotion.sh/docs/@emotion/core).

One of the biggest changes coming in this new version of Emotion is how the `css` prop works, both, on the API level and under the hood. In this post, I am going into detail on how the `css` prop works in Emotion 10.

## The `css` prop in previous versions of Emotion

Like the standard `styles` prop supported by React, Emotion’s  `css` prop has always allowed you to apply styles directly to your components or elements. For example, the following code would produce a paragraph with red text containing a `span` with blue text.

```js
import { css } from 'react-emotion'

const MyComponent = () => (
  <p css={{ color: 'red' }}>
    I have red text. And some{' '}
    <span
      css={css`
        color: blue;
      `}
    >
      blue
    </span>{' '}
    text.
  </p>
)
```

Notice how we are using object and string styles. The string styles use Emotion’s `css` function as a tagged template literal.

In order for this syntax to work, previous versions of Emotion required you to use the [Emotion Babel plugin](#). The plugin transpiles the `css` prop to a `className` prop whose value is a unique Emotion-generated CSS class name — a string. And all CSS classes for your component are part of a stylesheet that gets injected into the DOM at runtime.

Mitchell mentions some of the problems of the old implementation  in the Emotion 10 blog post.

> * It required a babel plugin  
> * It was not compatible with spreading an object as props  
> * Style composition order was confusing, undocumented and could break  

On top of these, one problem I always faced was the lack of support for theming. If you used the `emotion-theming` package to pass down a theme via the `ThemeProvider` component, you had to wrap your component in the `withTheme` [higher-order component (HOC)](https://reactjs.org/docs/higher-order-components.html), to use the theme in your styles.

```js
import { css } from 'react-emotion'
import { withTheme } from 'emotion-theming'

const MyComponent = ({ theme }) => (
  <p css={{ color: theme.colors.brandRed }}>
    I have red text. And some{' '}
    <span
      css={css`
        color: theme.colors.brandBlue;
      `}
    >
      blue
    </span>{' '}
    text.
  </p>
)

export default withTheme(MyComponent);
```

Emotion 10 addresses all of the above issues.

## Enabling the `css` prop in Emotion 10
You now have two ways of enabling the `css` prop in your React code. The first is to use the Babel plugin, just like in previous versions. However, this can be problematic in contexts where you do not want to or are unable to set up a custom Babel configuration. For example, you might be using [create-react-app](https://facebook.github.io/create-react-app/) or [CodeSandbox](https://codesandbox.io).

In those cases you can now resort to using Babel’s [`jsx` pragma](https://babeljs.io/docs/en/babel-plugin-transform-react-jsx#pragma). Use what? When I first read this in the Emotion 10 docs, I did not understand at all what this meant and when or how I should use it. So here is what I learned since.

Babel’s React JSX plugin has an option to specify the function it uses to transpile JSX expressions. The plugin defaults to using `React.createElement`. But by relying on `React.createElement`, you can only use syntax, such as props, supported by `React.createElement`. So, if you provide a different function to transpile JSX, you can support a different syntax.💡

Here is how you do this with Emotion 10.

```js
/** @jsx jsx */

import { css, jsx } from '@emotion/core'

const MyComponent = () => (
  <p css={{ color: 'red' }}>
    I have red text. And some{' '}
    <span
      css={css`
        color: blue;
      `}
    >
      blue
    </span>{' '}
    text.
  </p>
);
```

In the above example, we use the `jsx` pragma (`/** @jsx jsx */`) to tell Babel it should use the `jsx` function imported from `@emotion/core` to transpile JSX. And that is all you need to do to use the `css` prop. You do not need Emotion’s Babel plugin.

### When to use the `jsx` pragma

Initially, I got somewhat confused about when I should use import the `jsx` function and when to import `React`. ESLint would complain about unused `jsx` or `React` imports. So should I use `jsx` everywhere for the sake of consistency? The answer to this is a simple: “No”!

Importing React in a component is not an actual module import. It merely tells Babel that your JavaScript module contains JSX, which should be transpired using the React JSX plugin (i.e., with `React.createElement`). Using the `jsx` pragma allows you to opt into the `css` prop on a per component basis, by changing the way that component’s file is transpiled by Babel. You only need to specify the pragma, if you use the `css` prop in your component’s JSX. Otherwise, just keep importing `React`. You will definitely keep importing `React`, if your component uses something like `React.memo`, `React.Fragment`, `React.useState`, etc.

If you are using the Emotion Babel plugin, you do not need to use the pragma at all.

## How to use the `css` prop
You can still use the `css` prop as in previous versions of Emotion. The example from above will still work.

```js
/** @jsx jsx */

import { css, jsx } from '@emotion/core’;

const MyComponent = () => (
  <p css={{ color: 'red' }}>
    I have red text. And some{' '}
    <span
      css={css`
        color: blue;
      `}
    >
      blue
    </span>{' '}
    text.
  </p>
);
```

The styles will be evaluated and stored in Emotion’s cache under a class name. The class name is passed down to the React element through the `className` property. So far, so good.

One great improvement over previous versions is that you now have direct access to the theme context created by a `ThemeProvider`!

```js
/** @jsx jsx */

import { css, jsx } from '@emotion/core'
import { ThemeProvider } from 'emotion-theming'

const theme = {
  colors: {
    brandRed: 'red',
    brandBlue: 'blue',
  },
}

const MyComponent = () => (
  <p css={theme => ({ color: theme.colors.brandRed })}>
    I have red text. And some{' '}
    <span
      css={theme =>
        css`
          color: ${theme.colors.brandBlue};
        `
      }
    >
      blue
    </span>{' '}
    text.
  </p>
)

const App = () => (
  <ThemeProvider theme={theme}>
    <MyComponent />
  </ThemeProvider>
)
```

This is extremely powerful. Yes, some people might argue that having all these “inline” styles does not help with readability. I would tend to agree. However, styling an element inline every now and then, without having to wrap your component in a HOC or having to create a [styled-component](https://emotion.sh/docs/styled), *is* extremely useful.

## Conclusion
For me, the `css` prop has become a whole different animal with Emotion 10. Being able to use it without a Babel plugin is great when putting together a quick example in CodeSandbox or working with create-react-app. Having access to your application’s theme is invaluable when making smaller adjustments to existing components. Style composition with the `css` prop has also improved a lot with Emotion 10. The documentation on that is pretty great, so make  sure to [take a look](https://emotion.sh/docs/css-prop#style-precedence).