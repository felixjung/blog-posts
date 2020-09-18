# Presentational React Components — the 👩‍🎤emotion Way

I recently started working on a lot of mostly presentational React components in a project at work. These components are created with the [👩‍🎤 emotion](https://emotion.sh) [`styled` function](https://emotion.sh/docs/styled). `styled` mirrors more or less how [styled-components](https://www.styled-components.com) and [glamorous](https://glamorous.rocks) work. Instead of writing out a functional React component and adding styling to it, like this,

```js
import React from 'react'
import classNames from 'classnames'

// We wrote some BEM classes in a CSS file.
import './Input.css'

const Input = ({ disabled, isInvalid, ...props }) => {
  const className = classNames('input', {
    'input--disabled': disabled,
    'input--invalid': isInvalid,
  })
  return <input {...{ disabled, className, ...props }} />
}

export default Input
```

you end up writing this

```js
import styled, { css } from 'react-emotion'

const modifierDisabled = ({ disabled }) =>
  disabled &&
  css`
    pointer-events: none;
    opacity: 0.4;
  `

const modifierInvalid = ({ isInvalid }) =>
  isInvalid &&
  css`
    pointer-events: none;
    opacity: 0.5;
    border-color: red;
    color: lightcoral;
  `

const StyledInput = styled('input')`
  border-radius: 3px;
  border: 1px solid #999;
  font-size: 16px;
  padding: 8px 12px;

  &:focus,
  &:active {
    outline: none;
  }

  ${modifierDisabled} ${modifierInvalid};
`

export default StyledInput
```

Notice how pretty much all of this file has become styles? If you boil the code down to the "actual" JavaScript, you get this:

```js
import styled, { css } from 'react-emotion'

const modifierDisabled = ({ disabled }) => disabled && css``
const modifierInvalid = ({ isInvalid }) => isInvalid && css``

const StyledInput = styled('input')`
  ${modifierDisabled}
  ${modifierInvalid}
`
export default StyledInput
```

Personally, I find this very refreshing. If your component does not need logic, apart from "dynamic styles" (more on that later), you end up writing next to no code!

- You don't import your styles.
- You don't have to write a function for your component.
- You don't need to determine the classes you want to apply to your component with something like the classnames module.

 It feels very productive.  But what exactly is going on here?

## The `styled` and `css` Functions

The `styled` function is what creates your component. You pass in the component's style either via a tagged template literal or as arguments to a function call.

```js
// Tagged template literal usage
const StyledInput = styled('input')`
  color: black;
`

// Function usage with object styles
const StyledInput = styled('input')({ color: 'black' })
```

Emotion will create a class for your component and inject the stylesheet into the DOM.

In the examples further up, you've seen the [`css`](https://emotion.sh/docs/css) function. This function, like [`styled`](https://emotion.sh/docs/css), can also be used with tagged template literals and [object styles](https://emotion.sh/docs/object-styles). It can also take an array of objects, composing them into a single stylesheet. The function returns a class name, which you can then use either in a regular React component (i.e., without `styled`) or inside another `styled` or  `css` call. The latter is what we refer to as the _composition pattern_ in Emotion.

```js
const fontStyles = css`
  font-family: -apple-system,BlinkMacSystemFont,"Segoe UI", Roboto;
  font-size: 12px;
  line-height: 1.2;
`

const colorStyles = css`
  background-color: blue;
  color: yellow;
`
// As tagged template literal
const MyDiv = styled('div')`
  ${fontStyles}
  ${colorStyles}
`
// In a function call
const MyDiv = styled('div')(fontStyles, colorStyles)
```

Emotion will create three class names here. One for `fontStyles` , one for `colorStyles`, and finally another third class name for the `MyDiv` component. For the styles of `MyDiv` emotion will merge `fontStyles` and `colorStyles` and also take care of specificity for us.

This is a very powerful pattern, allowing you to create and re-use styles throughout your components. We will see how to leverage composition even more to create modifier styles for our components in the next section.

## Dynamic Styles

In the first example of using `styled` above, you see how style functions are passed to the `StyledInput` component. These functions receive an object parameter, which is destructured into variables. The variables are then used to determine whether to return styles or not. For example,

```js
const modifierDisabled = ({ disabled }) =>
  disabled &&
  css`
    pointer-events: none;
    opacity: 0.4;
  `
const StyledInput = styled('input')`
  // Base styles
  ${modifierDisabled}
`
```

What is happening here, is that Emotion calls every style function passed to a styled component, passing in the component's props. The style function then uses the props to calculate *dynamic styles*. This way we can compose multiple dynamic modifier styles of our component into one rendered style.

## Conclusion

Emotion's styled components — or those from styled-components and glamorous — enable simple yet powerful creation of presentational React components. This is facilitated by *dynamic styles* and *composition*.

In this post, I have used the way I've come to write styled components, with

- a styled component that has base styles
- and modifier functions that call `css`.

However, Emotion is flexible and very permissive in how you can use it. You may continue to write full BEM compliant stylesheets or put all your dynamic styles into the main `styled` call. This [issue on the Emotion repository](https://github.com/emotion-js/emotion/issues/381) discusses different approaches. Emotion is also compatible with a lot of tools from the large styled-components ecosystem. You can, for example, use the [polished "mixin" library](https://polished.js.org). It also offers a bunch of ecosystem packages of its own. There is a very refreshing approach to  media queries, called [facepaint](https://emotion.sh/docs/media-queries#facepaint). You can use the [ThemeProvider pattern](https://emotion.sh/docs/theming). There are [testing utilities for Jest](https://emotion.sh/docs/jest-emotion), that help snapshotting your styles. And so much more. Take a look at the [very nice documentation](https://emotion.sh/docs) and dive in (it lets you *edit examples in the docs*)! Finally, here is a [CodeSandbox with the input example](https://codesandbox.io/s/6v6y9qxw1w) from this post. Check it out.