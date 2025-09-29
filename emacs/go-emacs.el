;;; go-emacs.el -*- lexical-binding: t; -*-
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
;; Emacs org files to generate a static website. Instead of using the built-in emacs web publishing,
;; org-html-export-to-html to generate blog post snippets, which are then handed to go to build
;; the rest of the content.
;;
;;; Code:


(defun go-emacs-publish-blog ())

(setq go-emacs-blog-root-dir "/home/naokotani/code/go/go-emacs/")

(defun go-emacs-create-post ())

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

(provide 'go-emacs)
;;; go-emacs.el ends here
