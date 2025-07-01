---
title: "Notes for: Take My Money: Accepting Payments on the Web"
source_title: "Take My Money: Accepting Payments on the Web"
source_author: "Noel Rappin"
source_published: "2017"
source_edition: 1
ISBN: "978-1680501995"
categories:
  - notes
tags:
  - payments
toc: true
classes: wide
published: true
---

> [!info]
> title: {{ page.source_title }}
> author: {{ page.source_author }}
> published: {{ page.source_published }}
> edition: {{ page.source_edition }}
> ISBN: {{ page.ISBN }}

# Introduction

- this book is a simple introduction to collecting payments on the internet
- it shows examples of integrating to Stripe and Paypal API from a merchant's perspective
- it also shows some sample application architecture that a merchant can consider when accepting payments online
- I am focusing on the main design takeaways from the book.

# Architecture Design Notes

- **avoid floating point numbers** as these are imprecise. 
	- Whenever possible, we should be using integers to represent money related data.
	- Make use of industry standard libraries to handle the money type.
- **start with 3 layer architecture** containing controller, workflow, model (and database). 
	- This is simple to start with for a small merchant 
	- It can be refactored in the future to accommodate more complex business requirements
- **use a different workflow for a different provider** to help organise the code.
	- for example, have a specific workflow for Stripe integration and another one for Paypal integration
- **use Ngrok to enable developer callback tunnelling** 
	- for testing against the payment gateways (PG) to receive callbacks
- **gracefully handle failures**
	- validate bad data on client side before making requests to the PG
	- validate stale data on server side before making requests to PG (due to price/inventory changes while end user is checking out the cart)
	- throw readable errors for unhandled scenarios (and notify admin)
	- show actionable messages to end user
	- ensure data integrity even when payment fails (DB transaction rollback or compensating action)
	- switch to async processing to allow payment retry against PG without holding end user on session. Notify end user about eventual payment results (usually via email)
- **create admin access**
	- different roles of your staff requires different access level
	- make purchase on behalf of end user, which can bypass certain validation limitations (this is often very helpful for customer service)
- **model refunds as a new record**
	- so that payments that have already been completed are immutable
- **model discount and taxes with full replay-ability** 
	- so that an old payment can be correctly calculated, even when discount and tax formula changes in the future.
- **audit trail for data changes** so that we know which staff made what changes to the data