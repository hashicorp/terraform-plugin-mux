# Terraform Documentation

This directory contains the portions of [the Terraform website][terraform.io] that pertain to the Terraform Plugin Mux.

The files in this directory are intended to be used in conjunction with
[the `terraform-website` repository](https://github.com/hashicorp/terraform-website), which brings all of the
different documentation sources together and contains the scripts for testing and building the site as
a whole.

## Updating Sidebar Navigation

Updates to the sidebar navigation of Terraform docs can be made in the [nav-data] file. For details on structuring of sidebar navigation, see https://github.com/hashicorp/terraform-website#editing-navigation-sidebars

## Previewing Changes

You should preview your changes locally to ensure that the content is rendering properly before you create a pull request. The build includes content from this repository and the [`terraform-website`](https://github.com/hashicorp/terraform-website/) repository, allowing you to preview the entire Terraform documentation site.

To preview your content, complete the following steps:

**Set Up Local Environment**

1. [Install Docker](https://docs.docker.com/get-docker/).
1. Restart your terminal or command line session.

**Launch Site Locally**

1. Navigate into your local `terraform-plugin-mux` top-level directory and run `make website`.
1. Open `http://localhost:3000` in your web browser. While the preview is running, you can edit pages and Next.js will automatically rebuild them.
1. When you're done with the preview, press `ctrl-C` in your terminal to stop the server.

## Deployment

The website reads content from release tags to generate documentation for all versions of `terraform-plugin-mux` documentation. Changes merged into `main` will be included in the documentation for the next product release.

You cannot edit documentation for past versions of `terraform-plugin-mux` on the site. Documentation is an artifact of a product release. We push docs fixes forward for the next release, rather than retroactively fixing older versions.

[nav-data]: ../website/data/plugin-mux-nav-data.json
[terraform.io]: https://www.terraform.io/