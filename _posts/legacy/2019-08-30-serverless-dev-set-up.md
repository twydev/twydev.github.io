---
title: "Setting up a Serverless + NodeJS project"
toc: true
toc_label: "Sections"
published: false
---

How to set up a serverless nodeJS project from scratch. Includes instructions for Git, GitLab CI and Runner set up, ESLint, Jest, Winston, Complexity Report, and dotenv.

## Git Repo Set Up

1. create repository in GitLab
2. ensure SSH key has been generated and in place on local machine
3. clone repository to local machine
   ```
   git clone <ssh-url-of-repo>
   ```
4. configure local git user name and email (for this repo)
   ```
   git config user.name "<name>"
   git config user.email "<email>"
   ```
5. branch a new feature branch in the local machine (we are using GitHub flow)
   ```
   git checkout -b <feature-name-of-branch>
   ```
6. set upstream of this new branch
   ```
   git push --set-upstream origin <feature-name-of-branch>
   ```

## Start the Project (Serverless + NodeJS)

1. install yarn
2. create project
   ```
   yarn init
   ```
3. install serverless
   ```
   yarn add serverless
   ```

**Note** _running the following serverless commands in shell only worked because I had previously installed this package globally. Please use node to run the require serverless initialization or create npm scripts to do so_

4. create serverless project boilerplate using templates
   ```
   serverless create --template aws-nodejs
   ```
5. install serverless offline for local development
   ```
   yarn add --dev serverless-offline
   ```
6. add serverless plugin to serverless.yml file
   ```
   plugins:
   - serverless-offline
   ```
7. test the skeleton function locally
   ```
   serverless invoke local -f <functionName>
   ```

## Set Up GitLab Runner with AWS EC2 (Docker) and set up CI/CD

Read Hacker Noon [blog post](https://hackernoon.com/configuring-gitlab-ci-on-aws-ec2-using-docker-7c359d513a46) along side the documentations.

- [Runner Registration](https://docs.gitlab.com/runner/register/)
- [Runner Installation](https://docs.gitlab.com/runner/install/)
- [Runner configuration TOML file](https://docs.gitlab.com/runner/configuration/advanced-configuration.html)
- [Docker basic concept and usage](https://docs.docker.com/get-started/)
- [Find the necessary and ideally official Docker Images](https://hub.docker.com)
- [Another high level guide to CI/CD on GitLab](https://medium.com/@tarekbecker/a-production-grade-ci-cd-pipeline-for-serverless-applications-888668bcfe04)

1. Create EC2 instance (free tier / spot instance). Remember to get the SSH key.
2. Maximise the free tier eligible SSD storage (30 GB currently)
3. Assuming instance is running ubuntu, SSH to the instance,
   ```
   sudo apt-get update
   curl -l <gitlab's repository for Debian OS> | sudo bash
   sudo apt-get install gitlab-runner
   sudo gitlab-runner register
   ```
4. During registration, enter the gitlab-ci coordinator URL and registration token. These can be found in GitLab console, **Settings > CI/CD > Runner**
5. Set a tag for the runner, for example `private-ec2`. This allows us to specify the tagged runners to run specific jobs later.
6. Select Executor type: `docker`
7. Select docker image: `node:10.16.3`. As of this writing, this is the recommended LTS version.

This completes the registration. The runner should be working now. You can check GitLab console to see the runner.

**Important:** set the runner config in GitLab console to allow it to pick up all jobs that are untagged. Also, you might want to disable shared runner for the project to cut cost.

Next we need to install Docker.

1. Install required packages to use repository over HTTPS
   ```
   sudo apt-get install \
   ca-certificates \
   curl \
   software-properties-common
   ```
2. Add Docker's GPG key.
   ```
   curl -fsFL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
   ```
3. Setup the repository.
   ```
   sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
   ```
4. Update the package and install Docker.
   ```
   sudo apt-get update
   sudo apt-get install docker-ce
   # verify installation
   sudo docker run hello-world
   ```

Some other helpful commands:

```
sudo gitlab-runner list #list all runners
sudo gitlab-runner verify
sudo gitlab-runner stop
sudo gitlab-runner start
# since this installation made used to root user, the gitlab-runner config can be found in the following directory and you may need to modify some configurations. Remember to restart the runner when doing so.
/etc/gitlab-runner/config.toml
```

Next we need to set up a CI/CD config file in our git repository. This is the simplest example for our current project just to test that the setup is working.

1. Create a new `.gitlab-ci.yml` file in our repository root. (If we create this file in the repository through GitLab console, we can actually select CI templates pre-created for different software stacks)
2. Input the following contents.
3. Validate that the YAML file is in the correct format using GitLab console tool at the web address `gitlab.com/<project>/<repository>/-/ci/lint`.
4. Commit the code and push to the repository.
5. Check the pipeline status in GitLab console.

```
image: node:10.16.3

cache:
  paths:
    - node_modules/
    - .yarn

before_script:
  - apt-get update -qq && apt-get install

stages:
  - test
  - build

Test:
  stage: test
  tags:
    - private-ec2
  before_script:
    - yarn config set cache-folder .yarn
    - yarn install
  script:
    - echo "Successfully Ran Test on GitLab Runner"

Build:
  stage: build
  tags:
    - private-ec2
  before_script:
    - yarn config set cache-folder .yarn
    - yarn install
  script:
    - echo ""Successfully Ran Build on GitLab Runner"
```

## Set up Dev Dependencies (Code Style + Lint)

1. Install Prettier and ESLint. I opted out of editor config since prettier + linter does the job across all editors.
   ```
   yarn add prettier --dev
   yarn add eslint --dev
   # add 2 more dependencies,
   # eslint-config-prettier to deconflict formating responsibilities between Prettier and ESLint,
   # eslint-plugin-prettier to make ESLint run Prettier.
   yarn add -dev eslint-config-prettier eslint-plugin-prettier
   ```
2. Initialize ESLint.

   ```
   # for Linux
   eslint --init

   # for Windows
   node node_modules\eslint\bin\eslint.js --init
   ```

3. Add recommended config to .eslintrc
   ```
   "plugins": [
       "prettier"
   ],
   "extends": [
       "prettier",
       "eslint:recommended",
       "plugin:prettier/recommended"
   ],
   ```
4. Create prettierrc config file and specify the prettier rules
5. Create .prettierignore file so that all the rest of the file types do not get auto formatted.
   ```
   *.json
   *.yml
   *.md
   node_modules
   .eslintrc.js
   .prettierrc.js
   jest.config.js
   ```
6. For Sublime Text users, install JsPrettier plugin to auto format files on save. **Preferences > Package Settings > JsPrettier > Settings - User**.
   ```
   {
     "auto_format_on_save": true
   }
   ```
7. Configure **package.json** to run eslint. ESLint will now use Prettier for style checks due to the the set up in step3.
   ```
   "scripts": {
      "lint": "eslint ."
   },
   ```
8. Test to see if the linting works `npm run lint`.

## Set Up Unit Test Framework (Jest)

1. Install Jest
   ```
   yarn add --dev jest
   ```
2. Initialize Jest, which creates a config file

   ```
   # for Linux
   jest --init

   # for Windows
   node node_modules\eslint\bin\jest.js --init
   ```

3. Turn on code coverage in Jest config file.
   ```
   coverageDirectory: 'coverage',
   ```
4. Create a simple test and see if the test goes well `npm run test`

## Set Up Logger

1. I chose to use winston.
   ```
   yarn add winston
   ```
2. Set up a logger module. I have chosen to use a factory method to create the logger.

- the _metaMessage_ allows us to inject any additional default logging fields
- modules are only loaded once in the application, hence we will always be using a single logger instance.

```javascript
// logger.js
const winston = require("winston");

const logger = () => {
  const proto = {
    metaMessage: {},

    setMeta(message) {
      this.metaMessage = message;
    },

    info(message) {
      this.internalLogger.info(message, this.metaMessage);
    },

    error(message) {
      this.internalLogger.error(message, this.metaMessage);
    },

    warn(message) {
      this.internalLogger.warn(message, this.metaMessage);
    },

    debug(message) {
      this.internalLogger.debug(message, this.metaMessage);
    },

    internalLogger: winston.createLogger({
      format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.json(),
        winston.format.prettyPrint()
      ),
      transports: [new winston.transports.Console()]
    })
  };

  return Object.assign(Object.create(proto));
};

module.exports = logger();
```

## Set Up Complexity Report

1. Install complexity report library. This can help us track how complex our code is.
   ```
   yarn add --dev complexity-report
   ```
2. Set up a config file for complexity-report. Please refer to the npm page for more information https://www.npmjs.com/package/complexity-report
   ```
   //.complexrc
   {
     "output": "./.complexity/report.md",
     "format": "markdown",
     "allfiles": false,
     "ignoreerrors": true,
     "filepattern": "\\.js$",
     "silent": false,
     "newmi": true
   }
   ```
3. Create npm script and test run the library.
   ```
   "scripts": {
     "report": "cr ./src"
   }
   ```

Some metrics to be aware of (refer to https://radon.readthedocs.io/en/latest/intro.html for more information):

- **Cyclomatic Complexity**: number of decisions a block of code contains plus 1.
- **Cyclomatic Complexity Density**: ratio of Cyclomatic Complexity to SLOC.
- **Source Lines of Code (SLOC/LOC)**: number of lines of text in source code.
- **Halstead Complexity Measure**: Uses number of distinct operators and operands, and total number of operators and operands in the code to measure complexity.
- **Maintainability**: calculated using a factored formula consisting of Cyclomatic Complexity, SLOC, and Halstead Volume. (Microsoft Variant of the index is between 0 to 100)
- **Dependency Count**: number of CommonJS/AMD dependencies for the module.

## Setting up Environment Variables

Main goals of setting up environment variables:

- to work across different machines, but not commit sensitive information into the repository
- different values for different environment (dev, staging, prod)
- able to build the code smoothly in the local machine as well as during CI/CD

I chose to use dotenv and adopt a pattern of loading the environment variables through a config module. Articles for reference:

- [NodeJS Environment Variables](https://medium.com/the-node-js-collection/making-your-node-js-work-everywhere-with-environment-variables-2da8cdf6e786)

1. Install dotenv.
  ```
  yarn add dotenv --dev
  ```
2. Create a `.env` file with the environment variable.
  ```
  # .env file
  ACCESS_KEY=12345
  SECRET_KEY=12345
  ```
3. Create a config module to load the variables for your server
  ```javascript
  // config.js
  module.exports ={
    accessKey: process.env.ACCESS_KEY,
    secretKey: process.env.SECRET_KEY
  }
  ```
4. Access the variables in any modules simply by importing from config.js
5. To start up NodeJS or run Jest using `.env` file to provide the variables, we need to launch the server using the following option:
  ```
  // package.json
  scripts: {
    "start_with_env": "node -r dotenv/config server.js",
    "test_with_env": "jest --setupFiles dotenv/config"
  }
  ```
6. With Serverless Framework, the same set of variables from `.env` needs to be mapped to `serverless.yml` as `environment` properties of the function. If we are using GitLab premium for CI/CD, we can define different values for environment variables, and GitLab will expose the correct set of values to our runners depending on the executing environment (dev, stage, prod etc.). However, without premium, we may opt to identify our variables using names with environment as prefix (e.g. DEV_SECRET_KEY, PROD_SECRET_KEY). It will be up to our `serverless.yml` configurations to detect the current executing environment, and expose the correct variables to our function.
