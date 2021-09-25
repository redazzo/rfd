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
{{.RFD_first_state}}.

Update the record.md document in the RFD root directory accordingly.

### 4. Push Your RFD Branch Remotely

Push your changes to your RFD branch in the RFD repo.

    $ git add rfd/0042/README.md
    $ git commit -m '0004: Adding placeholder for RFD <Title>'
    $ git push origin 0004

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
