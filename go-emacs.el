;;; go-emacs.el --- Helper functions for go-emacs static site generator -*- lexical-binding: t; -*-
;;
;; Copyright (C) 2025
;;
;; Author:  <naokotani@naokotani>
;; Maintainer:  <naokotani@naokotani>
;; Created: September 28, 2025
;; Modified: September 28, 2025
;; Version: 0.0.1
;; Keywords: abbrev bib c calendar comm convenience data docs emulations extensions faces files frames games hardware help hypermedia i18n internal languages lisp local maint mail matching mouse multimedia news outlines processes terminals tex text tools unix vc wp
;; Homepage: https://github.com/naokotani/go-emacs
;; Package-Requires: ((emacs "27.2"))
;;
;; This file is not part of GNU Emacs.
;;
;;; Commentary:
;; This package is designed to work with go-emacs, a Go based static site generator that utlizes
;; Emacs org files to generate a static website. Instead of using the built-in Emacs web publishing,
;; org-html-export-to-html to generate blog post snippets, which are then handed to go to build
;; the rest of the content.
;;
;;; Code:

(require 'xdg)
(require 'org)
(defvar go-emacs-blog-dir
  (expand-file-name "go-emacs/" (xdg-user-dir "DOCUMENTS"))
  "Root directory of the go-emacs blog project.")

(defvar go-emacs-post-dir
  (expand-file-name "posts/" go-emacs-blog-dir)
  "Output directory for the go-emacs blog site.")

(defvar go-emacs-page-dir
  (expand-file-name "pages/" go-emacs-blog-dir)
  "Output directory for the go-emacs blog site.")

(defvar go-emacs-resume-dir
  (expand-file-name "resume/" go-emacs-blog-dir)
  "Output directory for the go-emacs blog site.")

(defvar go-emacs-output-dir
  (expand-file-name "blog/" (xdg-user-dir "DOCUMENTS"))
  "Output directory for go-emacs blog.")

(defvar go-emacs-package-dir
  (file-name-directory (locate-library "go-emacs"))
  "Directory where go-emacs.el is located.")

(defvar go-emacs-config
  ""
  "Path to the go-emacs config.toml file.")

(defun go-emacs-build ()
  "Build the go-emacs binary."
  (interactive)
  (let ((default-directory go-emacs-package-dir))
    (async-shell-command "make build" "*Build Go Emacs")))

(defun go-emacs-binary ()
  "Get the go-emacs binary location."
  (expand-file-name "go-emacs" go-emacs-package-dir))

(defun go-emacs-serve ()
  "Run the blog with Python web server."
  (interactive)
  (let ((default-directory go-emacs-output-dir))
    (async-shell-command "python3 -m http.server 8080" "*Go Emacs Serve*")))

(defun go-emacs-refresh-paths (root-path)
  "Refesh all directory paths based on `ROOT-PATH'."
  (interactive)
  (setq go-emacs-blog-dir   (expand-file-name root-path)
        go-emacs-post-dir   (expand-file-name "posts/" go-emacs-blog-dir)
        go-emacs-page-dir   (expand-file-name "pages/" go-emacs-blog-dir)
        go-emacs-resume-dir (expand-file-name "resume/" go-emacs-blog-dir)))

(defun go-emacs-publish-blog ()
  "Run the go-emacs binary with `async-shell-command'."
  (interactive)
  (let ((default-directory go-emacs-package-dir))
    (async-shell-command
     (concat (go-emacs-binary)
             " -d " go-emacs-blog-dir
             " -p " go-emacs-package-dir)
     "*Go Blog*"))
  (message "Publishing blog..."))

(defun go-emacs-get-parent-dir ()
  "Gets the parent directory of the current directory."
  (let ((dir default-directory))
    (with-temp-buffer
      (cd dir)
      (cd "..")
      (go-emacs-get-directory-name default-directory))))

(defun go-emacs-get-directory-name (dir)
  "Return name of DIR as string."
  (file-name-nondirectory (directory-file-name (file-name-directory dir))))

(defun go-emacs-publish ()
  "Publish smartly based on current directory."
  (interactive)
  (cond
   ((string= (go-emacs-get-parent-dir) (go-emacs-get-directory-name go-emacs-post-dir))
    (progn
      (message "Publishing post...")
      (go-emacs-publish-post)))
   ((string= (go-emacs-get-parent-dir) (go-emacs-get-directory-name go-emacs-page-dir))
    (progn
      (message "Publishing page...")
      (go-emacs-publish-page)))
   ((string= (go-emacs-get-parent-dir) (go-emacs-get-directory-name go-emacs-resume-dir))
    (progn
      (message "Publishing resume...")
      (go-emacs-publish-resume)))
   ((message "Not currently in a publish directory"))))

(defun go-emacs-create-post (dirname)
  "Create a new blog post in posts/DIRNAME with an .org file and images/ subdir."
  (interactive "sPost directory name: ")
  (let* ((posts-dir (expand-file-name go-emacs-post-dir))
         (post-dir (expand-file-name dirname posts-dir))
         (images-dir (expand-file-name "images" post-dir))
         (org-file (expand-file-name (concat dirname ".org") post-dir))
         (timestamp (format-time-string "[%Y-%m-%d %a %H:%M]")))
    (unless (file-exists-p posts-dir)
      (make-directory posts-dir t))
    (make-directory post-dir t)
    (make-directory images-dir t)
    (unless (file-exists-p org-file)
      (with-temp-file org-file
        (insert (format "#+title:\n#+date: %s\n#+filetags: :post:\n#+tags:\n#+summary:\n"
                        timestamp))))
    (find-file org-file)))

(defun go-emacs-create-page (dirname)
  "Create a new blog post in posts/DIRNAME with an .org file and images/ subdir."
  (interactive "sPage directory name: ")
  (let* ((pages-dir  (expand-file-name go-emacs-page-dir))
         (page-dir (expand-file-name dirname pages-dir))
         (images-dir (expand-file-name "images" page-dir))
         (org-file (expand-file-name (concat dirname ".org") page-dir))
         (timestamp (format-time-string "[%Y-%m-%d %a %H:%M]")))
    (unless (file-exists-p pages-dir)
      (make-directory pages-dir t))
    (make-directory page-dir t)
    (make-directory images-dir t)
    (unless (file-exists-p org-file)
      (with-temp-file org-file
        (insert (format "#+title:\n#+date: %s\n#+filetags: :page:\n"
                        timestamp))))
    (find-file org-file)))

(defun go-emacs-publish-post ()
  "Publish an org post into HTML."
  (interactive)
  (setq org-html-doctype "html5"
        org-html-html5-fancy t
        org-export-with-toc nil
        org-export-with-section-numbers nil)
  (org-html-export-to-html nil nil nil t)
  (go-emacs-publish-post-metadata))

(defun go-emacs-publish-post-metadata ()
  "Collecs metadata for a post and writes it to a metadata.toml file."
  (let* ((keywords (org-collect-keywords '("TITLE" "DATE" "TAGS" "SUMMARY")))
         (title   (car (alist-get "TITLE" keywords nil nil #'equal)))
         (date    (car (alist-get "DATE" keywords nil nil #'equal)))
         (tags    (car (alist-get "TAGS" keywords nil nil #'equal)))
         (summary (car (alist-get "SUMMARY" keywords nil nil #'equal))))
    (with-temp-file "metadata.toml"
      (insert (format "title=\"%s\"\n" title))
      (insert (format "tagString=\"%s\"\n" tags))
      (insert (format "summary=\"%s\"\n" summary))
      (insert (format "datestring=\"%s\"\n" date)))))

(defun go-emacs-publish-page ()
  "Publish an org page into HTML."
  (interactive)
  (setq org-html-doctype "html5"
        org-html-html5-fancy t
        org-export-with-toc nil
        org-export-with-section-numbers nil)
  (org-html-export-to-html nil nil nil t)
  (go-emacs-publish-page-metadata))

(defun go-emacs-publish-page-metadata ()
  "Collect metadata for a page."
  (let* ((keywords (org-collect-keywords '("TITLE" "DATE" "TAGS" "SUMMARY")))
         (title   (car (alist-get "TITLE" keywords nil nil #'equal))))
    (with-temp-file "metadata.toml"
      (insert (format "title=\"%s\"\n" title)))))

(defun go-emacs-publish-resume ()
  "Publish an org resume into HTML."
  (interactive)
  (setq org-html-doctype "html5"
        org-html-html5-fancy t
        org-export-with-toc nil
        org-export-with-section-numbers nil)
  (org-html-export-to-html nil nil nil t))

(provide 'go-emacs)
;;; go-emacs.el ends here
