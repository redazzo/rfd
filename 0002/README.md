---
title: Example RFD
authors: Gerry Kessell-Haak <gerry.kessellhaak@edpay.nz>
state: discussion
discussion: <link to discussion>
---

# RFD0002: Example RFD

The purpose of the Request for Discussion (RFD) process is to facilitate a lightweight means for raising, discussing, and accepting (or rejecting) ideas, concepts, designs, and decisions.

An RFD begins as a markdown document with a metadata header (as per this first RFD, as an example). The data to be captured includes the authors of the RFD, the state, and title. The state indicates where along the process the RFD has progressed, as per the following table. 

**NOTE:** *Never at anytime during the process do you push directly to the master branch. Once the pull request (PR) with the RFD in your branch is merged into master, then the RFD will appear in the master branch.*

| State | Description |
|--------|-------------|
|**ideation**|The sketch of an idea, a first cut, and may be discarded if the author decides to take it no further. Used to capture the beginnings of a thought, or even just a single sentence so that it's not forgotten. A document in the ideation state contains at least a description of the topic that the RFD will cover, providing an indication of the scope of the eventual RFD. An RFD in this state is effectively a placeholder |
| **prediscussion** | A document in the prediscussion state indicates that the work is not yet ready for discussion. The prediscussion state signifies that work iterations are being done quickly on the RFD in its branch in order to advance the RFD to the discussion state.|
|**discussion**| Documents under active discussion should be in the discussion state. At this point a discussion is being had for the RFD in a Pull Request.|
|**accepted**|Once (or if) discussion has converged and the Pull Request is ready to be merged, it should be updated to the accepted state before merge. Note that just because something is in the accepted state does not mean that it cannot be updated and corrected. |
|**implementing**| Not just published, work is actively progressing on implementing the idea or concept. |
|**committed**| Once an idea has been entirely implemented, it should be in the committed state. Comments on RFDs in the committed state should generally be raised as issues -- but if the comment represents a call for a significant divergence from or extension to committed functionality, a new RFD may be called for; as in all things, use your best judgment.|
|**abandoned**|If an idea is found to be non-viable (that is, deliberately never implemented) or if an RFD should be otherwise indicated that it should be ignored, it can be moved into the abandoned state. |

## The RDF Lifecycle

### Reserve an RFD Number

### Create a Placeholder RFD
