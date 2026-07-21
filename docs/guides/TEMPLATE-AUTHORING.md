# Template Authoring Guide

A template in Host Anything is a blueprint for deploying a service. It's written in TOML format.

## Step-by-Step

1. **Copy the Base:** Start by copying `templates/_base/template.toml`.
2. **Meta:** Fill out the `[meta]` section with the service name and version.
3. **Variables:** Define `[config]` variables for anything the user should customize (e.g., ports, passwords). Use best practices (secure defaults, descriptive text).
4. **Runtime:** Specify `type = "docker"` and the `image`.
5. **Storage and Ports:** Map `[volumes]` and expose `[network]` ports. Use variable substitution like `"${http_port}"`.

## Testing

Load your template into the local Host Anything instance by placing it in `/etc/hostanything/templates/`.

## Submitting to Marketplace

To submit a template:
1. Create a GitHub repo named `hostanything-template-<name>`.
2. Add the `hostanything-template` topic to your repository.
3. Include your `template.toml` and a `README.md`.
