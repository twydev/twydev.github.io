---
title: "Basic understanding of Web Browsers"
toc: true
toc_label: "Chapters"
---

Since the invention of JavaScript was solely driven by the desire to enrich web browsing experience through client side processing, we must know how modern browsers work, in order to truly know JavaScript. Here is a link to an [article written by Tali Garsiel](http://taligarsiel.com/Projects/howbrowserswork1.htm) that explores the internals of web browsers. The article is old, but gold.

**How Browers Work**

My takeaways from the article:

- HTML error tolerance comments from Webkit are hilarious.
- I am thankful for the efforts by our community of great engineers to create a robust set of standards that made our modern web browsing experience possible.
- It was nice to see concrete use cases of data structures like trees and stacks, as well as algorithms like bubble sort and merge sort, in an example that is relatable to everyone, the modern web browser.
- HTML document is first parsed into a DOM tree by main thread, with parallel threads fetching network resources and scripts. Scripts will block parsing and are executed synchronously.
- Styles are parsed into a separate tree.
- A render tree is constructed to mark the correct rendering order for each DOM node, and to associate the node with the corresponding computed style (elements of the tree are called renderer).
- The render tree is traversed to layout all elements on the viewport and to paint them (methods of the renderer object itself).
- Re-layout or re-painting can be applied to only a small subset of nodes.
- Changes to any renderers that requires re-layout or re-paint will fire off the corresponding events, triggering the layout or paint execution by the main thread.

**Web API Exposed to JavaScript**

We add functionalities to our web application by writing JavaScript that interacts with Web API. Interactions typically start with registering event handler to listen out for events. When the event gets emitted, the handler will execute our specified logic, and may update the rendering through another set of web APIs such as the DOM.

- Web API https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Client-side_web_APIs/Introduction
- Event References https://developer.mozilla.org/en-US/docs/Web/Events
