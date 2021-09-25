# Requests for Discussion

## Introduction

Almost all engineering of note is a social enterprise. A significant contributor to the rate of innovation and general productivity of an engineering team is the rate at which new ideas can be created, discussed, critiqued and improved, and either rejected or (hopefully) accepted and used.

The purpose of a Request for Discussion (RFD) process is to facilitate a lightweight means for raising, discussing, and accepting (or rejecting) ideas, concepts, designs, and decisions. The process and tooling as described here is based on the original W3C RFC process, and *heavily* influenced by [this blogpost](https://oxide.computer/blog/rfd-1-requests-for-discussion).

The following are examples of when an RFD is appropriate, these are intended to be broad:

* An architectural or design decision for hardware or software
* Change to an API or command-line tool used by customers
* Change to an internal API or tool
* Change to an internal process
* A design for testing

To facilitate the the above, we've created a simple commandline tool that integrates with Git to automate the initialisation of a repository, creation of an RFD, and management of an RFD process.

## Structure of an RFD and the RFD Repository

An RFD itself begins as a markdown document with a metadata header. The data to be captured includes the authors of the RFD, the state, title, and a link to the thread discussing the RFD. The state indicates where along the process the RFD has progressed, as per the following table. An example is shown below:

    ---
    title: Introduce Request for Discussion Process
    authors: Bob the Builder <bob@thebuilder.co>
    state: discussion
    discussion: <link to discussion>
    ---

An RFD repository is structured as follows:

1. Each RFD results in a folder named as per RFD ID, e.g. "0003". The folder holds a readme.md markdown document, and any other content, that describes the RFD. By using the readme.md naming convention tools and services such as github will automatically display its content.
2. A branch is created and named as per the RFD ID (i.e. in keeping the with the prior example, the branch name is "0003").
3. Discussion on pull a request becomes a record of the changes and updates.
4. Markdown metadata located at the head of the RFD document is used to capture the title, author(s), and other data, and used to populate an index of all RFDs.


![pic](./media/rfdrepo.png)

When a repository is created, the very first RFD is the RFD process itself, captured in a folder named 0001. A readme.md document in the 0001 folder is copied to the root of the RFD repository, and contains content similar that that found in this document, but targeted at the users of the RFD repository itself.

It is assumed that the process will be used as-is. However, by creating this directory we have provided the ability for others to update their own respective process by (if they wish) creating a "0001" branch and using the RFD process accordingly to capture the associated discussion.

## The RDF Process and Lifecycle

*Never at anytime during the process do you push directly to the master branch. Once the pull request (PR) with the RFD in your branch is merged into master, then the RFD will appear in the master branch.*

An RFD progresses through stages that default to the following. Note that these are configurable via the config.yml file (described later).

| State | Description |
|--------|-------------|
|draft|The first draft, can be used to capture the beginnings of a thought, or even just a single sentence so that it's not forgotten. A document in the draft state contains at least a description of the topic that the RFD will cover, providing an indication of the scope of the eventual RFD.|
|discussion|Documents under active discussion should be in the discussion state, with the discussion taking place in an active Pull Request.|
|accepted|Once (or if) discussion has converged and the Pull Request is ready to be merged, it should be updated to the accepted state before being merged into master. Note that just because something is in the accepted state does not mean that it cannot be updated and corrected.|
|committed|Once an idea is being acted on (e.g. being built, coded, or moved into an operational state), it is moved to the committed state. Comments on RFDs in the committed state should generally be raised as issues -- but if the comment represents a call for a significant divergence from or extension to committed functionality, a new RFD may be called for; as in all things, use your best judgment.|
|abandoned|If an idea is found to be non-viable (that is, deliberately never implemented after having been accepted) it can be moved into the abandoned state.|

### Prerequisites
The following assumes you have installed the RFD commandline tool. If you haven't, [instructions can be found here]().

### Initialising an RFD Repository

Ideally Github or similar is used to hold the RFD repository as you'll be able to take advantage of the default repository rendering and management tools.

The first requirement is to create the repository - [you can find instructions for Github here](https://docs.github.com/en/get-started/quickstart/create-a-repo).

Clone your newly-created repository, [**making sure you use ssh**](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository). This is more secure, and is required for the RFD commandline tool to function.

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

Fill in the RFD number and title placeholders in the metadata section of the new doc and add your name as an author. The status of the RFD at this point should be "ideation" or "prediscussion".

Update the record.md document in the RFD root directory accordingly.

### 4. Push Your RFD Branch Remotely

Push your changes to your RFD branch in the RFD repo.

    $ git add rfd/0042/README.md
    $ git commit -m '0004: Adding placeholder for RFD <Title>'
    $ git push origin 0004

*FUTURE STATE: The desired behaviour is that after your branch is pushed, the table in the README on the master branch will update automatically with the new RFD (IN PROGRESS!). If you ever change the name of the RFD in the future, the table will update as well. Whenever information about the state of the RFD changes, this updates the table as well. The single source of truth for information about the RFD comes from the RFD in the branch until it is merged.*

### 4. Iterate on the RFD in Your Branch
You can work on writing your RFD in your branch:

    $ git checkout 0004

Gather your thoughts and get your RFD to a state where you would like to get feedback and discuss with others. It's recommended to push your branch remotely to make sure the changes you make stay in sync with the remote in case your local gets damaged.

It is up to you as to whether you would like to squash all your commits down to one before opening up for feedback, or if you would like to keep the commit history for the sake of history.

### 5. Discuss Your RFD!
The beauty of this process is that we take advantage of GitOps-style pull requests.

When you are ready to get feedback on your RFD, make sure all your local changes are pushed to the remote branch. Change the status of the RFD from ideation or prediscussion to discussion then commit and complete a push:

    $ git commit -am '0004: Add RFD for <Title>'
    $ git push origin 0004

Once pushed, *open a pull request to merge your branch into the master.* After the pull request is opened anyone subscribed to the repo will get a notification that you have opened a pull request and can read your RFD and give any feedback.

The comments you choose to accept from the discussion are up to you as the owner of the RFD, but you should remain empathetic in the way you engage in the discussion.

For those giving feedback on the pull request, be sure that all feedback is constructive. Put yourself in the other person's shoes and if the comment you are about to make is not something you would want someone commenting on an RFD of yours, then do not make the comment.

### 5. Accept (or abandon) the RFD
After there has been time for folks to leave comments, the RFD can be merged into master and changed from the discussion state to the accepted state. The timing is left to your discretion: you decide when to open the pull request, and you decide when to merge it - use your best judgment. RFDs shouldn't be merged if no one else has read or commented on it; if no one is reading your RFD, it's time to explicitly ask someone to give it a read!

Discussion can continue on published RFDs! The discussion: link in the metadata should be retained, allowing discussion to continue on the original pull request. If an issue merits more attention or a larger discussion of its own, an issue may be opened, with the synopsis directing the discussion.

Any discussion on an RFD can always continue on the original pull request to keep the sprawl to a minimum.

If you feel your comment post-merge requires a larger discussion, an issue may be opened on it -- but be sure to reflect the focus of the discussion in the issue synopsis (e.g., "RFD 42: add consideration of RISC-V"), and be sure to link back to the original PR in the issue description so that one may find one from the other.
