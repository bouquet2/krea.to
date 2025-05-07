In this small blog post, I will talk about my experience with Nix and How to start using Nix-darwin and home-manager. This is an opinion piece and may/may not have factual information, keep that in mind.

## What is Nix?
Nix is a declarative and reproducible package-manager/build system hybrid tool. It has its own language for expressions such as packages and system configuration.

## My history with Nix
I've first discovered the existence of Nix in 2022. Back then, the documentation for Flakes was basically non-existent, and I knew nothing about it. So I just installed NixOS to my laptop and used the default configuration.nix system to configure my system.

After getting used to it, I've discovered the existence of home-manager. It was impressive how much I could do in the confines of my user. So I've transferred most stuff (not all) to home-manager, and used it as such.

After a short while, I was annoyed with the issues with NixOS such as not being able to run regular dynamic-linked binaries, which was really important for the work I was trying to do back then, and decided to move on the whole Nix life.

While most of my friends encouraged me to try out Nix again, I was hesitant to try it again because it doesn't fit to my workflow. Well, that changed in 2024.

I have finally replaced my old laptop with a Macbook Pro (Nov 2024). This meant I had to use macOS for a good while until Asahi becomes to the new Macbook's. Until then, I decided to tinker with the whole package management situation in macOS. It was at this time that I discovered nix-darwin. It looked like the perfect amount of Nix that was sufficient in my workflow, so I've decided to try it. I've sticked with it since, and believe it is pretty sufficient for most people.

Now that the history is out of the way, here is how to set this whole thing up.

## Setting up Nix-darwin
Setting it up is pretty simple, actually. The [README for nix-darwin](https://github.com/LnL7/nix-darwin) is well-written and offers multiple approaches of installing nix-darwin. I have gone with the Flakes approach this time so I've followed that part of the README.

Since I wanted this to be an all-inclusive configuration that includes my entire dotfiles, I've also decided to add home-manager to the mix. This was also pretty straightforward, just add the [home-manager repository to the inputs](https://github.com/kreatoo/dots-mac/blob/main/flake.nix#L8)

After this, I basically followed the respective documentation to see which options I can apply through this setup. Turns out, you can do quite a lot with nix-darwin, including but not limited to;

* Setting system preferences
* Installing Nix and Homebrew packages
* Writing launchd services in a simple manner (Best one by far IMO, check https://github.com/kreatoo/dots-mac/blob/main/hosts/akiri/apps.nix#L24 for an example)

So that is really it for the setup I've done. Only other thing to mention is that Thanks to [@getchoo](https://github.com/getchoo) I also managed to automate the flake update process, by [Having a GitHub Action that creates a PR with a cron timer](https://github.com/kreatoo/dots-mac/tree/main/.github/workflows)

I'm pretty comfortable with this setup for now, but stuff may always change in the future. If you need my configuration, it is available [here](https://github.com/kreatoo/dots-mac).
