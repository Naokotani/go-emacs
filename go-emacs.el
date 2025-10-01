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

(defvar go-emacs-root-dir
  (expand-file-name "go-emacs/" (xdg-user-dir "DOCUMENTS"))
  "Root directory of the go-emacs blog project.")

(defvar go-emacs-post-dir
  (expand-file-name "posts/" go-emacs-root-dir)
  "Output directory for the go-emacs blog site.")

(defvar go-emacs-page-dir
  (expand-file-name "pages/" go-emacs-root-dir)
  "Output directory for the go-emacs blog site.")

(defvar go-emacs-resume-dir
  (expand-file-name "resume/" go-emacs-root-dir)
  "Output directory for the go-emacs blog site.")

(defvar go-emacs-binary
  (expand-file-name "go-emacs" go-emacs-root-dir)
  "Path to the go-emacs blog binary.")

(defvar go-emacs-config
  ""
  "Path to the go-emacs config.toml file.")

(defun go-emacs-refresh-paths (root-path)
  "Refesh all directory paths based on `ROOT-PATH'."
  (setq go-emacs-root-dir   (expand-file-name root-path)
        go-emacs-post-dir   (expand-file-name "posts/" go-emacs-root-dir)
        go-emacs-page-dir   (expand-file-name "pages/" go-emacs-root-dir)
        go-emacs-config   (expand-file-name "config.toml" go-emacs-root-dir)
        go-emacs-resume-dir (expand-file-name "resume/" go-emacs-root-dir)
        go-emacs-config (expand-file-name "resume/" go-emacs-root-dir)
        go-emacs-binary     (expand-file-name "go-emacs" go-emacs-root-dir)))

(defun go-emacs-publish-blog ()
  "Run the go-emacs binary asynchronously to publish the blog."
  (interactive)
  (let ((process-environment (cons (concat "CONFIG_PATH=" go-emacs-config) process-environment)))
    (start-process "go-emacs-publish-blog"
                   "*Go Emacs Blog*"
                   go-emacs-binary))
  (message "Publishing blog..."))

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
  (interactive)
  (let ((org-html-doctype "html5")
        (org-html-html5-fancy t)
        (org-export-with-toc nil)
        (org-export-with-section-numbers nil))
    (org-html-export-to-html nil nil nil t))
  (go-emacs-publish-post-metadata))

(defun go-emacs-publish-post-metadata ()
  (interactive)
  (let* ((keywords (org-collect-keywords '("TITLE" "DATE" "TAGS" "SUMMARY")))
         (title   (car (alist-get "TITLE" keywords nil nil #'string=)))
         (date    (car (alist-get "DATE" keywords nil nil #'string=)))
         (tags    (car (alist-get "TAGS" keywords nil nil #'string=)))
         (summary (car (alist-get "SUMMARY" keywords nil nil #'string=))))
    (with-temp-file "metadata.toml"
      (insert (format "title=\"%s\"\n" title))
      (insert (format "tagString=\"%s\"\n" tags))
      (insert (format "summary=\"%s\"\n" summary))
      (insert (format "datestring=\"%s\"\n" date)))))

(defun go-emacs-publish-page ()
  (interactive)
  (let ((org-html-doctype "html5")
        (org-html-html5-fancy t)
        (org-export-with-toc nil)
        (org-export-with-section-numbers nil))
    (org-html-export-to-html nil nil nil t))
  (go-emacs-publish-page-metadata))

(defun go-emacs-publish-page-metadata ()
  (interactive)
  (let* ((keywords (org-collect-keywords '("TITLE" "DATE" "TAGS" "SUMMARY")))
         (title   (car (alist-get "TITLE" keywords nil nil #'string=))))
    (with-temp-file "metadata.toml"
      (insert (format "title=\"%s\"\n" title)))))

(defun go-emacs-publish-resume ()
  (interactive)
  (let ((org-html-doctype "html5")
        (org-html-html5-fancy t)
        (org-export-with-toc nil)
        (org-export-with-section-numbers nil))
    (org-html-export-to-html nil nil nil t)))

(provide 'go-emacs)
;;; go-emacs.el ends here
