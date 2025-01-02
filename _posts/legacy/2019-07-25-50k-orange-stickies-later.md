---
title: "50,000 Orange Stickies Later"
toc: true
toc_label: Insights
published: false
---

This talk covers the essence of Event Storming: use of color-coded Post-it notes to model our domain, involving all stakeholders to craft our Ubiquitous Language, Process Model, and rapidly develop UI that are user-centric.

**50,000 Orange Stickies Later** - Alberto Brandolini, Explore DDD Conference (Denver), 2017
{: .notice--primary}

Event Storming feels like a good tool that I can use to facilitate the conversations between Product Managers and Developers as we try to refine our Sprint Backlogs. Very interesting to hear how Brandolini advocates an approach which seems casual and intentionally imprecise can lead to greater clarity for everyone in the team. Here are my takeaways from this talk.

## Model the Whole Business Line using Domain Events

**Events** are represented by orange stickies. Verbs in past tense. Inputs from Domain Experts.
{: .notice--warning}

**Annotations** are represented by red/purple stickies. These are notes alongside Events, that can represent blockers or concerns.
{: .notice--danger}

We aim to make the whole process visible, allow massive learning across silo boundaries, and get to a consensus around the **core problem**.

The outcome is a big picture that enables multiple storytelling, described using incremental notation and a language that all the different stakeholders to can understand.

There is no scope limitation. The bottleneck is in the picture. The core domain is in the picture. Total clarity.

## Exploring Big Picture for Discovery

A big picture allows the team to vote and come to a consensus on the most important, most challenging, and the most uncertain part of te domain to tackle and prioritize in their development.

**Incremental Notation** means that after every round of discussion, previous conceptions of a solution design to the problem may be superseded by a more appropriate design in the light of newer and more accurate insights.

**Fuzzy by Design** means that we intentionally use imprecise definitions for terms to allow everyone and everything to be included in the discussion and the model. And this can also trigger interesting conversations.

## Process Modeling
Zooming into a portion of the big picture, to have a deeper discovery and analysis of the particular process (may be a feature, user story, epic, and etc.).

A process has some pre-conditions and some command to trigger it, some things happening in the system (the flow which we want to model), before producing a set of events and read/view models as outcome.

**Color-coded Process Mapping**

{% include figure image_path="/assets/images/screenshots/event-storm-process-model-color.png" alt="" caption="Process modeling using Colors" %}

**Policies** are rules that determine how the system is supposed to react on an **Event** and then generate a **Command**.
- Policies may be explicit.
- Policies may be implicit, the way people are behavior in the business because that is the way it works on the ground, without written agreement.
- Automated handling of events, like listeners, handlers, sagas, and process managers, are all included within the spectrum of Policies.

## Every Step Challenges Business Value
Every step can create or destroy value for given users as you discover more opportunities, assumptions, inconsistencies, and etc.

## Software Design
The same approach can be used to explore software design.

{% include figure image_path="/assets/images/screenshots/event-storm-software-design-color.png" alt="" caption="Software Design using Colors" %}

When investigating **Aggregates** the focus should be on 
- state machine logic, 
- behavior and not data, 
- and postpone the naming of aggregate, as premature naming is jumping to conclusions. We should be clear on the logic before giving the logic an appropriate name.

## Ubiquitous Language
Arriving at the point where everyone use the same set of terms takes time, and that is perfectly fine. 

Discovering the model that perfectly describes that domain is great but it is not in every stakeholder's interest to understand it. It may be sufficient for your users to just think that the system works like magic.

## Read Models
Returned to the Presentation to be consumed by the application user. This should not be data-centric, but should be designed in a way that helps users make the next decision in the application.

## Product Owner is essential
The Product Owner plays the role of challenging the model, and also answering any challenges to the business choice. This brings reality into the conversation instead of having only software theories in the discussions.

## Final Note
People will not self-organize around systems they do not understand.
