---
title: The Test Organisation Request for Discussion Process
authors: Gerry Kessell-Haak
state: discussion
discussion: 
---

# RFD0001: The Test Organisation Request for Discussion Process

| Description | Link |
|---|---|
|Installation and configuration instructions| [Installation](0001/installation.md)|
| Index of RFDs | [RFD Index](index.md) |

</br></br>
The tooling and process are based on the Request for Discussion Process [described here](https://github.com/redazzo/rfd)

# Introduction

The purpose of the Request for Discussion (RFD) process is to facilitate a lightweight means for raising, discussing, and accepting (or rejecting) ideas, concepts, designs, and decisions.

The following are examples of when an RFD is appropriate, these are intended to be broad:

* An architectural or design decision for hardware or software
* Change to an API or command-line tool used by customers
* Change to an internal API or tool
* Change to an internal process
* A design for testing

This process has been *heavily* influenced by [this blogpost](https://oxide.computer/blog/rfd-1-requests-for-discussion) at Oxide, including outright copying of large swathes of text! We hope they don't mind.

An RFD begins as a markdown document with a metadata header (as per this first RFD, as an example). The data to be captured includes the authors of the RFD, the state, and title. The state indicates where along the process the RFD has progressed, as per the following table. An example is shown below:

    ---
    title: Request for Discussion Process
    authors: Joe Bloggs <joe.bloggs@example.com>
    state: implementing
    discussion: <link to discussion>
    ---

| State | Description |
|--------|-------------|
|draft|The first draft, can be used to capture the beginnings of a thought, or even just a single sentence so that it's not forgotten. A document in the draft state contains at least a description of the topic that the RFD will cover, providing an indication of the scope of the eventual RFD.|
|discussion|Documents under active discussion should be in the discussion state, with the discussion taking place in an active Pull Request.|
|accepted|Once (or if) discussion has converged and the Pull Request is ready to be merged, it should be updated to the accepted state before being merged into master. Note that just because something is in the accepted state does not mean that it cannot be updated and corrected.|
|committed|Once an idea is being acted on (e.g. being built, coded, or moved into an operational state), it is moved to the committed state. Comments on RFDs in the committed state should generally be raised as issues -- but if the comment represents a call for a significant divergence from or extension to committed functionality, a new RFD may be called for; as in all things, use your best judgment.|
|abandoned|If an idea is found to be non-viable (that is, deliberately never implemented after having been accepted) it can be moved into the abandoned state.|


## The RDF Lifecycle

*Never at anytime during the process do you push directly to the master branch. Once the pull request (PR) with the RFD in your branch is merged into master, then the RFD will appear in the master branch.*

### 1. Reserve an RFD Number
You will first need to reserve the number you wish to use for your RFC. This number should be the next available RFD number from looking at the current git branch -r output.

### 2. Create a Branch For Your RFD
create a new git branch, named after the RFD number you wish to reserve. This number should have leading zeros if less than 4 digits. Before creating the branch, verify that it does not already exist:

    $ git branch -r *0004

If you see a branch there (but not a corresponding sub-directory in rfd in master), it is possible that the RFD is currently being created; stop and check with co-workers before proceeding! Once you have verified that the branch doesn't exist, create it locally and switch to it:

    $ git checkout -b 0004

### 3. Create a Placeholder RFD
Create a placeholder RFD with the following commands (make sure you're in the RFD root directory):

    $ mkdir -p rfd/0004
    $ cp templates/rfd.md rfd/0004/README.md

Fill in the RFD number and title placeholders in the metadata section of the new doc and add your name as an author. The status of the RFD at this point should be
