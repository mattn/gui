A widgets package for go
========================

This is an experiment in creating a "widgets" package for go.  The
idea is to create a widget library utilizing functions that combine
together multiple widgets to form larger widgets.  This is sort of a
more "functional programming" (and hopefully more go-like) style than
conventional GUI libraries, which tend to be strongly oriented toward
inheritence.

I would like to make the widget library at least somewhat independent
of the rendering back end.  The current back end is based on
websockets and html 5 with a bit of javascript in between.  I don't
think it's a very nice design, but the point is to create a nice API
that isn't dependent on the specific back end.  Not that I've created
this yet, but it's my goal.  I'd like to eventually be able to support
a GTK version (or native windows) of the same application with (close
to?) no change in the source code of the application.

So far, this framework is highly incomplete.  All I have are buttons
and text.  If this README gets out of date, I may have more widget
types implemented.

The easiest way to run the example is to install gb, and then run it
in the root directory.  Then you can type `bin/example` and point your
browser to `http://localhost:12345`.


What I would like to achieve
============================

The idea is to create a widget library that is extensible in the
following senses:

 - Users can create new high-level widgets without modifying the
   widget library itself.  As much as possible of the widget library
   itself should actually be created with the user-level interface, to
   ensure that it is both robust and easy to use.
 
 - Implementation will not be exposed to users, so the underlying GUI
   can be changed (e.g. from web-based to GTK-based, or from one
   web-based approach to another) without any changes to high-level
   (e.g. user-created) widgets.

 - User-created widgets can achieve anything that can be achieved with
   native widgets.  E.g. if there is a native color-selector, then you
   should be able to create an equivalent color-selector on your own
   (for instance if you wanted to have a color selector that only
   allowed colors that would be visible against a particular background).
   
 - There is a minimum of boilerplate required to create a new
   high-level widget.

 - The API is go-like:  e.g. interfaces are used to distinguish
   between different sorts of objects, e.g. those that can be clicked
   and those that have text that can be set.

I'm not sure how to achieve these goals.  The use of interfaces is
pretty obvious, but how do we enable easy extension, and hiding of
implementation details? An obvious way to hide implementation details
is to put the implementation into a hidden method in the interface.
That leaves the question of how users will create new widgets.  By
hiding the implementation, users are forced to create new widgets
using exported functions or types.  But if we hide the data members of
our raw data types, it becomes hard for users to embed these in their
own types.
