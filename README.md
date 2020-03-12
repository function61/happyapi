[![Build Status](https://img.shields.io/travis/function61/onni.svg?style=for-the-badge)](https://travis-ci.org/function61/onni)
[![Download](https://img.shields.io/badge/Download-bintray%20latest-blue.svg?style=for-the-badge)](https://bintray.com/function61/dl/onni/_latestVersion#files)

REST API for delivering happiness - hosted on AWS Lambda.

tl;dr: put [this URL](https://function61.com/api/happy)
to your application to enable your users to get their daily dose of happiness.

NOTE: this is a very new project and the URL will change to a prettier domain soon.


Use case
--------

I wanted to have a "Enjoy your day!" wish at the footer of a web app I offer for my users.
I wanted the "enjoy" word to be a link that takes the user to a random picture on the
internet that brings happiness:

![](docs/example-ui.png)

**click** gets you:

![](docs/example-happiness.png)


Can I too use the URL?
----------------------

Yes! And don't be afraid to use it - I make the following promises:

- The URL is the API and it won't change, or if it will the old URL will get redirected (i.e. still work)

- The pictures will be family friendly

- The service won't have ads, or if in the long term will have ads they will be unobtrusive text-only ads.


Adding new pictures
-------------------

This currently only works for me - the project maintainer. Public submission may be coming later.

- Use [Online UUID Generator](https://www.uuidgenerator.net/) to generate UUID like this:
`10e239c4167f` (it's the two parts between both sides of the first dash).

- Upload images to `onni.function61.com` S3 bucket. Remember to set `Content-Type: image/jpeg` (or similar)

- Add references to main.go, remember to sort lines!


How to deploy
-------------

If for some reason you want to host your own API (you could just use the public API that
we host), follow these instructions.

Deployment is easiest using our [Deployer](https://github.com/function61/deployer) tool.
You don't need it and you can upload Lambda zip and configure API gateway manually if you want.

You have to do this only for the first time:

```
$ mkdir deployments
$ version="..." # find this from above Bintray link
$ deployer deployment-init happy-api "https://dl.bintray.com/function61/dl/onni/$version/deployerspec.zip"
Wrote /home/joonas/deployments/happy-api/user-config.json
```

Now edit above file with your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`

Then do the actual deployment:

```
$ version="..."
$ deployer deploy happy-api "https://dl.bintray.com/function61/dl/onni/$version/deployerspec.zip"
```
