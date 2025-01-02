---
title: "Basic understanding of NodeJS"
published: false
---

*Work in progress.*

**Why does my command not work in shell but works in npm run-script?**

I had installed serverless as a devDependencies of my project. I have also created a corresponding npm script.

```
// package.json
"scripts": {
	"local": "serverless invoke local -f hello -s dev"
}
```

In the root of this project, I ran the exact same bash command in npm script:

```
serverless invoke local -f hello -s dev
# /bin/bash: line 89: serverless: command not found

# However, running npm run-script works fine
npm run local
> serverless invoke local -f hello -s dev
```

**Package Installation**s

I have installed serverless package locally and as a devDependencies. When we run the shell command using npm run-script, we are able to launch everything within the project's node_modules.

> In addition to the shell’s pre-existing PATH, npm run adds node_modules/.bin to the PATH provided to scripts.
> Any binaries provided by locally-installed dependencies can be used without the node_modules/.bin prefix.

Source: [npmjs doc](https://docs.npmjs.com/cli/run-script)

And I had the misconception that I am always able to run serverless in shell with no problem on my previous machine, as long as `npm install` successfully downloaded all the packages into node_modules. But in reality, this was only possible because I had previously installed the serverless package globally for a project in the past.

> In npm 1.0, there are two ways to install things:
>
> 1. globally —- This drops modules in {prefix}/lib/node_modules, and puts executable files in {prefix}/bin, where {prefix} is usually something like /usr/local. It also installs man pages in {prefix}/share/man, if they’re supplied.
> 2. locally —- This installs your package in the current working directory. Node modules go in ./node_modules, executables go in ./node_modules/.bin/, and man pages aren’t installed at all.

Source: [nodejs global vs local install](https://nodejs.org/en/blog/npm/npm-1-0-global-vs-local-installation/)
