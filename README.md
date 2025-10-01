# Go Emacs
go-emacs is a static site generator that leverages the Emacs' `org-html-export-to-html` and Go html templating to produce a Static blog site. Includes `go-emacs.el`, which provide functions to create and publish posts and pages that will be included in the blog. It depends on a go binary, but is designed to be a provide a seamless blogging experience inside Emacs via org-mode.

# Example website
To see a sample website generated from these files, see [This Sample](https://naokotani.github.io/)

# Setup
go-emacs depends on a  Go binary, `go-emacs`, that can be build from the provided source code. The `Makefile` provides three commands `build`, `run` and `serve`. Build will build the binary from the provided source code and put it in the root directory. `run` will build and run the binary, creating an output directory and building the website. `serve` depends on both `build` and `run` and will create a simple testing web server using Python. The Python dependency is not required, and any web server that can server a static website from a directory will suffice.

# Emacs setup
The provided `go-emacs.el` file must be required in Emacs by adding the following to your `init.el`:

```
(add-to-list 'load-path "path/to/go-emacs")
(require 'go-emacs)
```

This will load the following interactive functions designed to be used with go-emacs:

## Functions
### Change Root Directory
Use `(go-emacs-refresh-paths "/path/to/project/root")` to change the directory that the helper functions point to. This will update location variables to reflect the newly defined root variable.

### Create Page
`(go-emacs-create-page)`
This will create a new "page" in the page directory. A page will appear on the top Nav, and will not contain post metadata such as the publish date. It will prompt for a directory name *which will also be used as the URL slug for the post*. There is not currently any mechanism to ensure these are unique. I considered adding the date string to the slug, but decided this would make for ugly URLs, and decided to leave it to the user to ensure unique post names.

### Create Post
`go-emacs-create-post)`
Similar to `(go-emacs-create-page)`, but creates a post. A Post will be displayed on the home page, and includes a time of publishing time stamp as well as `tags` designed to group themed posts.

### Publish Page
`(go-emacs-publish-page)`
This will 'publish' a page by producing the `.html` file required to generate the page. It will also generate a `metadata.toml` file which provides other information needed to generate the page such as the page title. 

### Publish Page
`(go-emacs-publish-post)`
Similar publish page, but produces the requisite `metadata.toml`. A post can optionally have a thumbnail that is optionally displayed on the post card if `cards=true` is set in the `config.toml`. The thumb should be in the same directory and be called `thumb.png`. During the build process, If the thumb width is greater than 400px, it will be resized to 200px to optimize the home page as it may have a larger number of images.

### Publish Resume
`(go-emacs-publish-resume)`
The resume page exists to optionally add an HTML resume to the blog along with a link to download a pdf.
This function will publish the currently opened page by producing an HTML version of the org file. I didn't create a special function this, but instead provide an example org file that can be replaced with a personal resume. To link a pdf, add the path to the resume in the `config.toml`. To skip the page entirely, set `isResume` in the `config.toml` to false and it will be skipped in the build process.

### Publish Blog
`(go-emacs-publish-blog)`
This is a convenience function that simply runs the `go-emacs` binary. Alternatively, you can run the binary on the command line as usual and it will work as long `CONFIG_PATH` variable is set to the location of a valid `config.toml`, or if the site files are in the default `~/Documents/go-emacs`

## variables
`go-emacs.el` defines a few variables that it uses to find to locations of your files.

### Root Dir
`go-emacs-blog-root-dir` There is no reason to change this directly as it is only used as a basis to locate the other required files and directories. If you want to change this, use `(go-emacs-refresh-paths)` function above.

### Location Directories
`go-emacs-post-dir`, `go-emacs-page-dir`, `go-emacs-resume-dir`, and `go-emacs-binary` can be changed individually with `setq` if they are not located in the project root directory. This must be done *after* calling `(go-emacs-refresh-paths)`, which will reset the path back to the root directory.

### Config Path
go-emacs uses a config.toml file to control many variables used in site generation, and the build process will fail entirely if the file is not provided. It is located in one of two ways; by checking for the value of the `CONFIG_PATH` environment variable, or by looking for it in `~/Documents/go-emacs/config.toml`. NOTE: This variable is not the directory, but the exact location of the file itself, eg. `/path/to/config.toml` NOT `/path/to/config/directory`. If the file is in the default location, this should remain set to an empty string, `""`.

# Config
The `config.toml` file is required to build the blog and is located in the root directory. By default, `go-emacs` will look for it in `~/Documents/go-emacs`, but the location can be changed by running `go-emacs` by setting the `CONFIG_PATH` variable to an absolute path to the `config.toml` file. Further documentation is provided within the file itself to explain the various settings, but at a minimum the file locations must be set if go-emacs is not located in `~/Documents/go-emacs`.

## Some Options of Note
The resume page is optional, and setting `isResume=true/false` to false will cause the build process to skip it when the site is built. The `pdf` option is blank by default and the 'download pdf' button won't be generated, putting in an absolute path to a pdf file will cause it to be copied to the resume folder when the site builds. There are two versions of the home page. A 'card' style and a more old school list style. By default, `cards=true/false` is set to true and the home page will generate in this cards style, but when set to false the list style will generated. Try 'em both! When the site builds  `rss=true/false` is set to true, meaning RSS icons and an RSS XML page will be generated. If you aren't interested in RSS, you can set it to false to skip building it.

# Styles
Included in the `static/css` directory are two files `normalize.css`, which is a CSS reset and `styles.css`, which provides the base styles for the blog. A third file `vars.css` is created when the website is built based on the values set in the `styles.toml` file. This provides a simple way to change the fonts, font sizes, and colours for the blog. Of course, further customization can be achieved by editing `styles.css` itself, which is recopied each time the website is built. 

# Post and Page Images
Both 'posts' and 'pages' can have linked images in them. To link an image, place the image in the `images/` directory for that post/page and then link the file normally in the org document with `file:/images/image.png`. When the site builds the images will be copied to an image folder for that post/page.
