---
title: "Lessons from my first product launch"
toc: true
toc_label: "Lessons"
---

I have participated in the launch of a B2B e-commerce platform. Here are some of the lessons I have learnt in terms of resource management, engineering process, as well as the preparation leading up to the launch.

## Outsource vs In-house Engineers

We were able to build the software from scratch and launch it in 6 months thanks to the help of a huge outsource team of engineers. Outsourcing increased our speed of going to market and allowed us to push out a lot more features than what our in-house team was capable of delivering.

However, we had to struggle with poor code quality, miscommunications and wrong assumptions due to language barrier, and high rate of bugs.

## Engineering processes can be improved

A number of things were not put in place before we started the project, and it resulted in massive technical debts just so that we can go to market on time. If we were to start a new product all over again, here is what I would do in sequence:

1. **Have a proper network plan on AWS.** Since we want to have CICD and infrastructure-as-code for our projects, we need to ensure that all foundational shared infrastructures are already put in place. Regardless of the tools used (TerraForm, CloudFormation), we should make sure that VPC, Subnets, RouteTables, ACLs, DNS / Route53 configurations are in place. The biggest challenge is how to draw a line between common resources and project specific resource.

2. **Adopt GitHub flow. (which can only work with CICD and Infrastructure-as-Code).** Instead of following the traditional Git flow (which was unfortunately what we have done), we should have used the GitHub flow to ensure that whatever feature we are adding always follows the configuration of Production environment. Unlike our current state, where Prod, Staging, and Dev environments had all deviated from each other, which makes debugging much more difficult. What we should have done is to

   - build new features in a new branch
   - pass tests in the branch
   - deploy branch to production (rollback if buggy)
   - merge into master and delete branch

3. **Test Driven Development.** After all the frustrations from buggy software, I believe it is only basic courtersy and good hygiene to cover your own code with unit tests. Especially if the component / feature you are tasked to deliver involves complex logic. One of the absolute worst practice is to rely on a QC to catch bugs for you, or to let you know that your code failed to meet the business requirements. Why would you expect another person to flush the toilet for you?

4. **Code Review.** I am not sure how code review can ever work in an outsource context, but I believe it should be done. I did not have enough time to ensure good code quality and I also did not want to be a bottleneck to the development, therefore all merge requests were easily passed. In the end, the software was is laced with anti-patterns, violation of SOLID principles, code smell, which all contribute to high defect rate.

Some of the less critical but important improvements that should be made:

- proper AWS IAM execution roles management (for system)
- proper AWS access management (for users)

## Preparation before Launch

The number one lesson is to not overwork. I pushed pretty hard, sleeping 4 hours a day for almost 2 weeks just to meet the deadline of launch. However, our launch date got delayed due to critical bugs still not fixed on time, and every additional day of delay means that I need to continue my unsustainable lifestyle for another day, until the software gets launched. Situation did not improve post launch due to other operational issues occuring. In the end, I was burning out and I could not function efficiently for the next 2 weeks. In conclusion, it was simply not worth it to go overtime.

> With every overtime comes an undertime to restore the balance.

2. **Data Migration Efficiency.** As pointed out by my senior colleague, it would have been much more efficient to copy out the entire database of the old system as temporary tables in the database of the new system, and use SQL scripts or functions to perform the data migration, instead of using scripts of various programing languages. This ensures no data loss, and we can re-run the migration any time with no latency.

3. **End-to-end Tests.** Due to the drift in configurations between environments, we could not carry out proper end-to-end tests across 3 systems, which includes the entire customer journey in the frontend. What happened was a critical bug occured in production post launch but this bug somehow eluded all tests, and could only be detected by a downstream system. This can be avoided with sufficient end-to-end testing.

4. **Rundown Checklist.** What I found useful was to have a rundown checklist, detailing what exactly needs to be done at what time, with key information (IP addresses, UUIDs, commands etc.). So that we can all be kept on track during the actual launch and also see the progress.

5. **Monitoring Dashboard, not just for system health.** We had dashboards set up, but it was only to monitor usage and detect errors. They did not turn out to be as helpful when critical bugs occured that denied our customer's orders from flowing to a downstream system, with absolutely no errors thrown.
   - What we needed was key business metrics (such as number of orders, summary table of order IDs etc.) at each checkpoints of our software platform, so that we can monitor whether an order has flown through from end-to-end.
   - Also, we could have designed some quick sanity checks on all the important data to ensure our software is working well. These checks can be scheduled to run periodically.
   - We can have automated checks on frontend experience at regular schedule, and after every software updates, to ensure that our user experience have not degraded.
   - We should have a daily checklist for our engineers to run through everyday, at least for the first 2 weeks of launch, to check if system is stable.
