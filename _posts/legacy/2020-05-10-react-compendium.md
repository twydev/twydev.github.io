---
title: "My React & React Native Compendium"
toc: true
toc_label: "Chapters"
published: false
---

React is a frontend library that is performant and lightweight. In order to use React well, we need to understand the basics of lifecycle, states, and props, the difference between React and React Native, and some basic understanding of Redux and Testing in React.

My notes are derived mostly from [reactjs.org](https://reactjs.org/) and [edX course, HarvardX: CS50M](https://courses.edx.org/courses/course-v1:HarvardX+CS50M+Mobile/).

- React is declarative (interested in _what_ to execute, not _how_ to execute).
- Paradigm encourages decomposing problem to be handled by smaller reusable components.
- React has performant Reconciliation process to sync app state to DOM.
  - reconstruct Virtual DOM.
  - find difference between Virtual DOM and actual DOM.
  - only apply necessary changes to actual DOM.

---

## JSX

JSX = JavaScript XML, is a language used by React that will be transpiled to JavaScript. React acknowledges that logic and UI are coupled, and it may be more beneficial to be able to declare both using a single langauge. Separation of concern can be acheived by organizing our code into reusable components.

- lowercase tags are HTML tags.
- uppercase tags are React component tags.
- ReactDOM escapes any values from JSX and convert them to string before rendering, protecting us against injection attacks.

## React Element

JSX objects are transpiled into React elements. Elements are the building blocks used to render to the DOM. These are pure objects, therefore they are cheaper to create and maintain, as compared to DOM elements.

- ReactDOM is responsible of rendering React elements. (Need to specify the target DOM node).
- React Elements are immutable. Re-rendering require us to call ReactDOM render method again with new Element.
- Even though a new Element instance is passed to ReactDOM, only the necesary changes will be reflected to DOM, making it more efficient to update.
- Developers can simply focus on what an entire element look like at every point in time, instead of worrying about how an update can be applied.

---

## React Components

React Components are reusable JavaScript functions that form the building blocks of React application logic.

- accepts an input object (commonly known as "Props", which stands for properties).
- returns a React element to be rendered.

If a component tag is used to to define a React Element, any JSX attributes and children to this tag will be passed as a single Props object to the component function.

_Creating a Component as a JavaScript class provides access to additional functionalities such as state and lifecycle. If we don't we will be creating Stateless Functional Components (SFC)_

### State and Lifecycle

As long as we render a component (class) into the same DOM node, only a single instance of the class will be used. This allows us to use `state` to manage our Component internals over time, instead of strictly depending on Props. The class also inherit a set of lifecycle methods that runs at certain critical juncture of React Component's lifecycle.

```JSX
class MyComponent extends React.Component { // need to extend React Component class
  constructor(props) {
    super(props); // must call parent constructor to initialize context object/THIS. ES6 class behavior
    this.state = { date: new Date() };
  }
  // this.props is accessible as it will be initialized by parent constructor.

  componentDidMount() { // Run once, after initial render of component
    this.timerID = setInterval( // add a timer as member of the class to retain reference
      () => this.tick(), // callback function triggers a class method, using lexical THIS through arrow function.
      1000
    );
  }

  componentWillUnmount() { // Run once, when component is getting removed from the DOM
    clearInterval(this.timerID);
  }

  tick() {
    this.setState( // calls setState to update component state. This triggers React re-rendering
      { date: new Date() }
    );
  }

  render() {
    return (
      <div>
        <h2>Time Now: {this.state.date.toLocaleTimeString()}</h2>
      </div>
    );
  }
}

ReactDOM.render(
  <MyComponent />,  document.getElementById('root')
);
```

_Example from [reactjs.org](https://reactjs.org/docs/state-and-lifecycle.html)_

State must be updated with `setState()` method.

- State updates may happen in batch, and is asynchronous.
- When using object as input to `setState()`, the object will be merged with current state. (Since updates may be batched, all the inputs will be merged, therefore we may not observe all updates rendered to DOM).
- Callback function can be used as input to `setState()`. Callback should expect two inputs `(state, prop)`, the first being the existing state before update, and second is the props at the time of update. Callback should return the updated state object.

### Props

Props are immutable, with respect to their React components.

`this.props` accesses `props` property, which is specially reserved (like `state`) and will be initialized by `React.Component` class, which all components will inherit from.

`props.children` will be populated if the current component has been called with child element included within the component tag in JSX.

#### PropTypes

- Only runs in **development mode**.
- Validate type of component props at runtime.
- Serves as documentation of component API.

To use, we need to import PropTypes, and attach types declaration to a React component. Since it provides additional checking, we should use it. However, if we are using TypeScript with React, this may be unnecessary.

### Data Flow

- Component state is encapsulated without component. Parent or child components should not know whether a component is stateful or stateless.
- We can allow state updates to flow down from parent to child, by passing it as props to the child.
- Therefore, data can only flow down, to the next component below in the tree.

### Losing Callback Context

As with all callback functions in JavaScript, if we pass a class method as a callback event handler, when it is eventually invoked, reference to context object `this` will be lost.

We can do a few things to get the desired behavior we want:

1. explicitly bind `this` from within the class to the method when registering the callback.
2. explicitly bind `this` to those methods within constructor.
3. use lexical `this` instead. Define class methods using arrow function, which statically inherit `this` context from the class (since the class is the enclosing scope of the arrow function).
4. use lexical `this`. Only pass arrow functions as callbacks. When registering the callbacks, define arrow functions that internally calls class methods, which inherits lexical `this` context.

This is an expected JavaScript behavior and is not a limitation of React or JSX.

[reactjs.org](https://reactjs.org/docs/handling-events.html) recommends the 2nd or 3rd approach to avoid extra re-rendering in child components. (If we define an arrow function in our props passed to a child component, the child will re-render whenever the parent element re-render and re-create that arrow function);

### Handling Events

- There is no need to explicitly call `addEventListener`, since we can simply register a handler when an element is rendered.
- Event default behavior can be suppressed using `event.preventDefault`
- If we want to pass extra parameters to our callback handler, we have two approach depending on how we are binding our context.

```JSX
<button onClick={(event) => this.deleteRow(extraParams, event)}>Delete Row</button>
<button onClick={this.deleteRow.bind(this, extraParams)}>Delete Row</button>
```

An event handler function can be set as callback at multiple locations in the React element, and it can vary its behavior based on e.g. `name` of the element that triggered the event, or other input properties.

### Render Null

Returning `null` in a component render method will prevent it from rendering, however the component will still go through the component lifecycle as usual (related lifecycle methods will still be called).

### Render Lists

JavaScript arrays can be first mapped to a list of `<li>...</li>` elements. An array of such element can be directly injected into JSX and will be rendered as a HTML list (provided we insert the appropriate `<ul>` tags).

React expect such list items to contain a key. `<li key={number/string}>`

This helps React to figure out which list item was updated and needs to be rendered. Therefore key should be derived from the item values and can uniquely identify the item. We can also resort to simply using index of item in the array as key (will incur huge performance penalty if ordering of items in the array is prone to changes).

The key can be specified when declaring a list item component as an element in JSX.

```JSX
function ListItem(props) {
  return <li>{props.value}</li>;
}

const listItems = numbers.map((number) => {
  <ListItem key={number.toString()} // this key do not get passed into the component, but will render the list properly in React. Therefore component code has no access to this key.
    value={number} />
);
```

### Controlled Components

Certain HTML elements has their own internal states and will conflict with React state if not managed properly.

The recommended approach is to let React state be the "single source of truth" and update those HTML element states using React states, creating _controlled components_.

The idea is straightforward:

- value of HTML element is solely based on React component state.
- events from HTML element will be processed by handlers within React component. (sometimes we may need to suppress event default behavior, to allow React component full control over that HTML element).
- these handlers extra the required update, and update React component state.
- this triggers re-rendering, and updates value of HTML element.

### Lifting State

{% include figure image_path="/assets/images/screenshots/react-lifting-state-up.png" alt="" caption="Lifting State Up" %}

_Visualisation of components used in this [example](https://reactjs.org/docs/lifting-state-up.html) on reactjs.org_

- Essentially, if a few children/sibling components need to share states, then their parent (common ancestor) will need to be the one maintaining that state.
- There can be no two way binding or sharing of state between components, since **data flow downwards**.
- Using IOC, we can allow child components to trigger state changes in the parent component, and these changes will propagate to all other child components via props.

---

## Design Pattern

React documentation guide recommends using Composition instead of Inheritance to achieve any software designs we want.

If we want to focus on business logic, then we should implement it in a separate JavaScript module, and import it into React component only as APIs.

React application should focus on what it does best, building good UI. Therefore, Composition should be sufficient.

### Thinking in React

Summary of an example thinking process provided by React documentation guide:

1. Create a mock UI
2. Break UI into Component Hierarchy
3. Build a static version of UI in React (sample data)
4. Identify the minimal state representation necessary to make UI interactive. (e.g. if we need both data and statistics of those data to provide interactivity, we should only represent the data in our state, and derive statistics only when we require them).
5. Identify where states should live (consider lifting state to common ancestor if necessary).
6. Add inverse data flow (egistering callbacks that introduce IOC for child components to trigger update on parent component states).

---

## React Native

### Introduction

- Based on React Core. JavaScript gets bundled for mobile. Write once run anywhere.
- Unlike browsers that runs rendering and JavaScript on one main thread, on mobile, separate threads will be used for UI, Layout, and JavaScript. In other words, JS thread maybe blocked but UI is still working.
- Communication between threads goes through a bridge and is asynchronous.
- Has different base components, elements, style, and navigation from React web.
- No access to browser API. (some functions have been poly filled)
- Not globally scoped, React Native needs to be imported.

### Style

- Uses JavaScript objects for styling instead of CSS file.
  - Object keys are CSS properties.
  - `StyleSheet` can be used to reuse styles and optimize sending style information over the bridge.
- Flexbox layout (default to column layout).
- Lengths are unitless.
- `style` props can take an array of styles.

### Event Handling

- Only a few "touchable" elements in React Native can accept event handler
- Callback function interface is not consistent. Need to consult documentation before use.

### Lists

#### ScrollView

Most basic scrolling view. (Unlike React Web, we will not be able to scroll unless we use a component with such interactivity).

- Components in an array need a unique key prop.
- Renders all children before appearing.

#### FlatList

Performant scrolling view that "lazily" renders components.

- Only currently visible rows will be rendered in the first cycle.
- Virtualized. Rows are recycled, and rows that exits visibility may be unmounted. (Note: component state will be lost if unmounted).
- Takes an array of data, and a `renderItem` function as props. (refer to docs for renderItem function signature).

#### SectionList

Extends `FlatList` with sections.

- Each section has its own data array and can override its own `renderItem` function
- `renderSectionHeader` function used to render header.

### User Input

React recommends using controlled components to render input, making react component states the source of truth.

- `onChange` prop is used to trigger react component when input is detected, to update states and re-render.
- `value` prop is used to pass the output state content to be rendered on input element.
- Input validation can be triggered by handler, before setting component state.
  - `setState` takes a second callback function as input, which will be called after setting state. This can be used for validation.
  - `componentDidUpdate` can also be used to call validation as well. (must compare current state with previous state to know when to NOT validate, to prevent causing infinite set state loop from the validation).
- Need to manually handle form submission (React Native limitation) and obtain all require values from component state.
- `KeyboardAvoidingView` is good for simple and short forms. It prevents virtual keyboard from obstructing the view.
  - Specify behavior when virtual keyboard appears by adjusting padding, height, or position of view.
  - However, the view moves independently from any child text inputs, so may not be good for complex form.

### Navigation

- Web navigation is oriented around URLs
- Mobile navigation API is completely different, on both iOS and Android.
- React Navigation library provides an platform agnostic alternative

#### Navigator, Routes and Screens

- Navigator is a React component that implements a navigation pattern.
- Route is a child of a navigator.
  - Each route must have a name and a screen component to be rendered when route is active.
  - The screen component can also be another navigator, creating a nested structure.
- The navigation prop is automatically passed to each screen component, allowing screen component to navigate to other route.
- `screenProps` can be used for rapid prototyping to pass prop to ALL screen component through navigator, but every route will re-render when screenProps change, so it is not efficient.
- When navigating from route to route, we can pass in a `param` to pass state to a different route.
  - We can't use props since that was pass through the parent navigator.
  - similar concept URL param in web.

#### Switch Navigator

- Inactive screens are unmounted.
- The only action allowed is to switch from one route to another.

#### Stack Navigator

- State of inactive screens are maintained and remains mounted.
- Platform specific layout, animations and gestures.
- Screens can be pushed/popped from the stack, or replaced.

#### Tab Navigator

- State of inactive screens are maintained and remains mounted.
- Provides a tab bar to switch between tabs.
- Platform specific layout, animations and gestures.
- By default `goBack()` returns to first tab.

#### Composing Navigator

- A navigator can be a screen component of a route under another navigator.
  - Do not render a navigator within screen component.
  - Instead, set the screen component as a navigator from the parent component.
- An app should only have one top level navigator, but we can navigate to any route in the app.
- `goBack` works for the whole app (supports Android back button).

### Development Tools

#### Expo

Expo is a suite of tools to accelerate React Native development. Comes with:

- Snack - run React Native in browser.
- Client - run project on mobile devices while developing.
- SDK - exposes cross-platform libraries and APIs.

#### Debugging

- React error/warnings
  - `console.error` triggers a full page error display
  - `console.warning` triggers a yellow banner (will not appear in production mode)
- Chrome Developer Tool (devtools)
  - JavaScript running inside Chrome tab is monitored by the debugger
  - Since React Native runs UI and JS on separate thread, communicating asynchronously through a bridge, they can be running on separate devices. When using chrome debugger, our JS app is running on chrome as a service worker.
  - Chrome debugger is able to log circular JSON, which other debugger may not support.
- React Native Inspector
  - similar to browser devtools inspector, able to click on layout elements on screen and look at element details.
  - but cannot live edit elements.
- `react-devtools` library
  - inspect component hierarchy, component props and state
  - install as project package dependency
  - run with CLI command (works with React Native Inspector)
  - allows live edit of style, props, etc.

#### Performance

- React Native Perf Monitor.
  - shows refresh rate of UI and JS threads.
- Chrome Performance Profiler.
  - only in development mode.
  - flame chart showing time taken to render each component.

### Common Performance Inefficiencies

- Rendering too often
  - Props changes not related to UI was passed to component.
  - If using Redux, component should only subscribe to interested state changes.
  - Use keys in array/list to prevent re-rendering the same element.
  - `shouldComponentUpdate()` lifecycle method can be used to provide you the control over whether a component should re-render.
  - `React.PureComponent` has a default `shouldComponentUpdate()` method that performs a shallow diff of props, to determine if it should re-render.
- Unnecessary props changes
  - Passing a new props object to a component may cause re-rendering of the entire sub-tree.
  - Any object literals or functions created in `render()` method will also cause new object to be created at each render (immutable)
  - Use constant, methods, or properties of class instance to avoid unnecessary immutable object.
- Unnecessary logic in mount/update (quite subtle)
  - Properties created in a class will be re-created at every mount.
  - However, class methods will only be created once ever (prototype chain).
  - The disadvantage of using class methods: we cannot use arrow functions, and hence we cannot rely on lexical scope of `this` context, and will instead need to perform explicit binding.
- Animation
  - Animation on UI thread that requires JS thread data to be sent over the bridge rapidly burdens performance.
  - Blocking on either threads greatly reduce UX.
  - Implementing animation in native (Swift, Java) may not be easy to manage project.
  - Use `Animated` API to declare computation in JS, to be run on UI thread instead. (but we can't use native driver for layout props).

## Redux

### Flux Architecture

Redux is inspired by the Flux architecture. A Flux architecture has the following characteristics:

- Unidirectional data flow.
- Views react to changes in a **store**.
- Only a **dispatcher** can update data in a store.
- Dispatcher can only be triggered by **actions**.
- Actions are triggered from **views**.

Motivation of using such an architecture is to manage complexity of numerous models and views with multidirectional dependencies in a large React application. Some advantages that Redux offers:

- direct management of deeply nested component state.
- reduce duplication of information in state.
- reduce risk of not updating props or not passing props properly.
- clarity of overall app state.

Redux data management uses a single store. **Actions** trigger **Reducers** which update the **Store**.

#### Reducers

- Pure functions. No side effects. Output determined by inputs.
- Takes previous state and an action to return new state.
- A new state object should be returned for immutability.

#### Store

- Expose getter functions to obtain current state.
- Can only be updated using dispatch function (which works with reducers).
- Allows adding of listeners to run callbacks when state changes.

#### Actions

- Data object containing information for state update.
- Usually contains a `type` key to indicate the type of state update.
- Created by Action Creator Functions.
- Actions must be dispatched.

### react-redux

- Official React bindings for Redux
- Recommended to use Higher Order Component manage state updates and re-rendering using this library.
  - `connect()` wraps a React component (e.g. a screen component)
  - Register callback `mapStateToProps` whenever state updates. Props will be passed to wrapped component.
  - Register callback `mapDispatchToProps` to bind action creator function to store dispatch, and expose it to wrapped component as props.
  - Library needs to register a store using `Provider`, usually at the top level component. The registered store will be available to all nested components that are using the `connect()` API.

### Redux Middleware

- Allows us to extend redux without modifying redux implementation.
- Any function with the following prototype can be a middleware:
  - `({getState, dispatch}) => next => action => void`
  - The middleware is a function that takes a store as an input and returns a second function.
  - The second function takes a `next` middleware as input, and returns the third function.
  - The third function takes an `action` as input, and execute without returning. (It has access to all inputs so far due to closure over function scope).
  - `next` param allows us to chain the middleware.
- Intercepts and modify incoming action to the reducer.
  - Can be used to trigger and respond to async calls before updating state.

### redux-persist

- abstracts the persistance of store (uses `AsyncStorage` under the hood).
- allow us to persist store when app closes and re-opens.
- display loading screen while waiting for store to rehydrate.

### Container vs Presentation Components

To manage complexity as our application grows, we can consider:

- Having components that are aware of redux state (containers).
- Components that only renders what was passed as props (presentation).

### Testing

#### Jest for Redux Actions

Besides using the standard Jest test functions, we can make use of snapshots.

- Test script can save snapshot of function output from first run, and compare it with subsequent run.
- Will throw error if snapshot does not match.
- We can choose to update snapshot with new outputs.

#### Jest for Async Redux Actions

- We may use mock functions or mock external libraries.
- Dependency injection pattern in our code design will help us pass mock functions.

#### Jest for React Component

- We can also use snapshots to compare render output.
- `react-test-renderer` allows rendering of a component outside the context of an app.
- `jest-expo` has config required for React Native testing.

---

## Additional Reading

- React Advanced Guide [page1](https://reactjs.org/docs/accessibility.html)
- Why props update also cause re-rendering, besides setState? Or is it because re-rendering of parent (passing new props down) also re-render child components?
- Full Component lifecycle?
- How to do testing properly, since React is built with test in mind.
- Props are not considered as updated and will not render if the reference are not updated? (e.g. sort array in place will not render. Create a new sorted array will render)
- Higher Order Component.
- Commonly used icon pack `react-native-vector-icons`
- smart and dumb components by [Dan Abramov](https://medium.com/@dan_abramov/smart-and-dumb-components-7ca2f9a7c7d0)
