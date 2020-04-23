.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Documentation

.. _rs_dev_docs:

Developing Documentation
========================

As an open ecosystem project, we encourage community feedback and involvement.  Docs can be updated by
pull requests against the `github <https://github.com/digitalrebar/provision>`_ repository either from a private
tree or directly against the tree.

A couple of notes about consistency.

#. Digital Rebar is the name of the parent project and can be abbreviated DR.
#. Digital Rebar Provision or DR Provision or DRP can be used to reference this part of the project.
#. API docs generated from the go files as part of swagger annotations of the godoc comments.  Update there, please.
#. CLI docs generated from the cli files as part of cobra structures.  The tools generate those.  Update there, please.

Otherwise, try and find a good place for what needs to be added.  And Thanks!


Documentation Tooling
---------------------

There are a lot of ways to work with ReStructuredText...  Below is only one possible way of setting up a working environment to ensure you are writing clean RST based documentation.  This method is designed to be as lightlweight as possible, while still being as accurate as possible with the final rendered doc changes.   There are a lot of editors tha will render RST formatted markup, but very few of them render it the same way, or similar enough to the final rendered documentation styles to be correct.

This process uses the following elements:

#. Install Sphinx to correctly render the RST markup text to HTML
#. Edit the doc in any text editor of your preference (eg *vim*, *atom*, *emacs*, etc.)
#. Use a terminal window with a ``watch`` function to rebuild the HTML document tree
#. Use a web browser with an auto-refresh extension to view the rendered HTML


Install and Setup Sphinx and Swagger
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

These setup steps were tested and verified on a Mac platform.  However, they should be the same for Linux.

1. Pre-requirement is to have ``pip`` installed:

  ::

      sudo easy_install pip


2. Install Sphinx:

  ::

    pip install Sphinx

.. note:: for more detailed information, please review the Sphinx website on how to install it, at:
  http://www.sphinx-doc.org/en/stable/tutorial.html#install-sphinx


3. Install SwaggerDoc Library

Perform this step from the root of the project

  ::

    sudo -H pip install -r requirements.txt

.. note:: If you receive an error message on your first HTML tree build (when you run ``make html``) similar to:
    ::

       sphinx-build -b html -d _build/doctrees   . _build/html
       Running Sphinx v1.6.5
       Extension error:
       Could not import extension sphinxcontrib.swaggerdoc (exception: No module named swaggerdoc)
       make: *** [html] Error 1


  You will need to build and install the *swaggerdoc* components manually.  To do so, do the following:
    ::

      git clone https://github.com/galthaus/sphinx-swaggerdoc
      cd sphinx-swaggerdoc/
      sudo python setup.py  install
      cd ..
      rm -rf spinx-swaggerdoc/


Edit the Docs
~~~~~~~~~~~~~

  * Checkout/clone the Digital Rebar Provision repo from Github
  * Modify the doc(s) as appropriate
  * Verify the modifications are rendered correctly and fix errors/warnings
  * Create a branch
  * Submit a pull request for your changes


Rebuild the HTML Rendered Docs
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To assist with editing/reviewing your changes, you will want to rebuild the RST format documentation to rendered HTMl.  The Sphinx package is used to do this.  Generally, you simply change to the *base git cloned Provision* directory, and do:

  ::

    make html

If you are making a lot of changes, you will want to use a file watching utility to see watch for file writes, and automatically kick the ``make html`` process off for you.

On Mac OS X - you can install the ``fswatch`` package:

  ::

    brew install fswatch

An example use would be:

  ::

      # cd to the base Provision git directory
      fswatch -o -0 -r doc | xargs -0 -n 1 -I {} make html

Now, any time any changes are made to the files in the ``doc/`` directory, the ``make html`` process will be automatically kicked off for you.  You should run this in a separate terminal window.


View Your Doc Edits Locally
~~~~~~~~~~~~~~~~~~~~~~~~~~~

Once you've run ``make html`` as above, the _RSTr_ format files will be built into _HTML_ generated versions.  You can view the effects of your changes by opening your web browser up, and navigating to the _HTML_ generated docs on your local disk.  The below path example references my home directory, and github path location.  You'll have to modify this to your local User and location where you've stored your github clone.

Point your browser to the on-disk rendered location (which is the ``_build/`` directory in the base git repo on disk).  For example:

  ``file:///Users/shane/github/digitalrebar/provision/_build/html/doc/dev/dev-docs.html``


Auto-Refresh Browser
~~~~~~~~~~~~~~~~~~~~

The last piece of the puzzle, you will want to set your web browser to auto-refresh a given tab or window.  This way, the HTML rendered documentation will be refreshed in the browser.   There are several add-ons/extensions that will do this for you.  Here at RackN we have used the following extensions:

  Chrome *Auto Refresh Plus* extension:
    https://chrome.google.com/webstore/detail/auto-refresh-plus/hgeljhfekpckiiplhkigfehkdpldcggm

  Firefox *Tab Reloader* add-on (works on Chrome, Firefox, and Opera; but limited to 10 second reloads as minimum reload time):
    https://add0n.com/tab-reloader.html

Simply set your browser tab to refresh every 5 or so seconds.

Final Steps Before Committing
-----------------------------

Once you are statisfied with your changes, you need to do a complete clean build of the doc tree.  To do this, you do the following:

  ::

    rm -rf _build
    make html

Fix all warnings and errors you introduced.  If you are authoring or fixing docs for content-packages or plugin_providers, remember to
use the `doc-override` directory to test your built content pack docs.


Hints and Tips for Content Packs and Plugin Providers
-----------------------------------------------------

Here are some tips for building and writing documentation for Content Packs and Plugin Providers.

Content Pack RST File
---------------------

For a content pack, you will need to do the following to get the documentation file from the content pack.  For this example, we will
assume that your content pack is in the directory, *example*.  You will need to do the following steps.  Only the last is different from
your probable normal test procedure.  This also assumes that `drpcli` is in your path.

  ::

    cd example
    drpcli contents bundle ../example.yaml
    drpcli contents document ../example.yaml > ../example.rst

At this point, you can copy the `../example.rst` file to the `doc-override` directory in your `digitalrebar/provision` tree and follow the same
build and view process.



Plugin Provider RST File
------------------------

For a plugin provider, you will need to use the `tools/build-one.sh` command.  Once you completed editing your content section of your Plugin Provider,
you will need to build it.  Using `example` again, you would do the following:

  ::

    tools/build-one.sh cmds/example

This will generate an `example.rst` in the `cmds/example` directory.  This file can then be copied to the `doc-override` directory in your
`digitalrebar/provision` tree and follow the same build and view process.


Header Section Levels
---------------------

The file ``._Documentation.meta``, inside a content pack or the content portion of a plugin provider, should be RST format.  The build tools will automatically
bundle the content pieces into a build product file.  This fill will be upload to an Amazon S3 bucket when the build completes.  The sphinx config file, ``conf.py``,
controls what gets included from the Amazon S3 bucket and downloaded in the ``content-packages`` directory.  The ``content-packages.rst`` file is a simple
all-inclusive TOC of files contained in ``content-packages``.

Within the ``._Documenation.meta`` file, section separations must follow this heirarchy because the tools add pieces to the top to make the page consolidate and
show in the table of contents correctly.

  ::

    ~~~~~~~~~~~ - Reserved for the Title of the content pack or plugin provider
    ----------- - Next level down - all new sections in ._Documenation.meta should at the level
    =========== - Next level down - within the higher sections
    +++++++++++ - Next level down - within the higher sections
    ^^^^^^^^^^^ - Next level down - within the higher sections

The goal of the ``._Documentation.meta`` insert is that it can add a descriptive set of information at the highest level and then start creating sub-sections as
needed.  The build process will append second level (``-------------``) sections for all the included object types within the content.


Here is an example of a ``._Documentation.meta`` file in the example content package:

  ::

    This is the main descriptive section.

    SubSection1
    -----------

    SubSection1Sub1
    ===============

    SubSection1Sub2
    ===============

    SubSection2
    -----------

    SubSection2Sub1
    ===============


If the content package, ``example``, were rendered it would produce a single file:

  ::

    .. Copyright (c) 2017 RackN Inc.
    .. Licensed under the Apache License, Version 2.0 (the "License");
    .. Digital Rebar Provision documentation under Digital Rebar master license
    .. index::
      pair: example; Content Packages

    .. _rs_cp_example:


    example
    ~~~~~~~

    This is the main descriptive section.

    SubSection1
    -----------

    SubSection1Sub1
    ===============

    SubSection1Sub2
    ===============

    SubSection2
    -----------

    SubSection2Sub1
    ===============

    params
    ------

    This content package provides the following params.

    example/cool-param
    ==================

    Documentation entry from the example-cool-param.yaml file.


    <<< for all the included object types >>>

The single file can be built by running, ``drpcli contents document example.yaml``.  The required input is
a content package bundle file.  This will generate an RST file to stdout.  Use the normal bundling process to
generate the yaml or json file.
