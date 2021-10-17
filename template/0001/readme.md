---
title: {{.Title}}
authors: {{.Authors}}
state: {{.State}}
discussion: {{.Link}}
---

# RFD0001: {{.Title}}

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
    id: 0003
    title: Request for Discussion Process
    authors: Joe Bloggs <joe.bloggs@example.com>
    state: implementing
    discussion: <link to discussion>
    ---

| State | Description |
|--------|-------------|
{{ $state := "" -}}
{{ $desc := "" -}}
{{ range $k, $v := $.RFDStates -}}
    {{- range $sk, $sv := $v -}}
        {{- range $ssk, $ssv := $sv -}}
            {{- if eq $ssk "name" -}}
                {{- $state = $ssv -}}
            {{- end -}}
            {{- if eq $ssk "description" -}}
                {{- $desc = $ssv -}}
            {{- end -}}
        {{- end -}}
|{{ $state }}|{{ $desc }}|
{{ end }}
{{- end }}

## The RDF Lifecycle

*Never at anytime during the process do you push directly to the master branch. Once the pull request (PR) with the RFD in your branch is merged into master, then the RFD will appear in the master branch.*

### 1. Creating an RFD

To create an rfd, simply issue the new command while in the root of the rfd repository i.e.

    $ rfd new

This will:
* Issue a new rfd id based on the highest of the branch id, the id of the local directories, the id of the remote directories, and the highest remote branch id.
* Request information on the title and authors of the rfd (this could potentially be automated by picking up the git user and email details).
* Create and check-out the local branch with the name as per the allocated id
* Create the rfd directory and readme.md file.
* Stage, commit, push, and set the upstream branch of the rfd

When done, you'll automatically be on the new branch. To edit (using the nano editor as an example):

    $ cd 0002
    $ nano readme.md

Of course, you can use whatever editor you're comfortable with.

### 2. Iterate on the RFD in Your Branch

Gather your thoughts and get your RFD to a state where you would like to get feedback and discuss with others. It's recommended to pull and push your branch remotely on a regular basis to make sure the changes you make stay in sync with the remote.

It is up to you as to whether you would like to squash all your commits down to one before opening up for feedback, or if you would like to keep the commit history for the sake of history.

### 3. Discuss Your RFD!
The beauty of this process is that we take advantage of GitOps-style pull requests, so everything is as per the normal git process.

When you are ready to get feedback on your RFD, make sure all your local changes are pushed to the remote branch. Change the status of the RFD to discussion then commit and complete a push:

    $ git commit -am '0002: Add RFD for <Title>'
    $ git push origin 0002

Once pushed, *open a pull request to merge your branch into the master.* After the pull request is opened anyone subscribed to the repo will get a notification that you have opened a pull request and can read your RFD and give any feedback.

The comments you choose to accept from the discussion are up to you as the owner of the RFD, but you should remain empathetic in the way you engage in the discussion.

For those giving feedback on the pull request, be sure that all feedback is constructive. Put yourself in the other person's shoes and if the comment you are about to make is not something you would want someone commenting on an RFD of yours, then do not make the comment.

### 4. Accept (or abandon) the RFD
After there has been time for others to leave comments, the RFD can be merged into master and changed from the discussion state to the accepted state. The timing is left to your discretion: you decide when to open the pull request, and you decide when to merge it - use your best judgment. RFDs shouldn't be merged if no one else has read or commented on it; if no one is reading your RFD, it's time to explicitly ask someone to give it a read!

Discussion can continue on published RFDs! The discussion: link in the metadata should be retained, allowing discussion to continue on the original pull request. If an issue merits more attention or a larger discussion of its own, an issue may be opened, with the synopsis directing the discussion.

Any discussion on an RFD can always continue on the original pull request to keep the sprawl to a minimum.

If you feel your comment post-merge requires a larger discussion, an issue may be opened on it -- but be sure to reflect the focus of the discussion in the issue synopsis (e.g., "RFD 42: add consideration of RISC-V"), and be sure to link back to the original PR in the issue description so that one may find one from the other.
