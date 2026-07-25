# gtr

Go Translate — a minimal CLI for Google Translate.

## Usage

```sh
gtr [flags] <text to translate>
```

Simplest form — auto-detect source language, translate to English:

```sh
$ gtr 猫
cat
```

```sh
$ gtr -t ja cat
猫

[noun]
	猫: cat

[noun]
	- a small domesticated carnivorous mammal with soft fur, a short snout, and retractable claws. It is widely kept as a pet or for catching mice, and many breeds have been developed.
	- short for catalytic converter.
		"models fitted with a cat as standard"
	- short for catamaran.
[verb]
	- raise (an anchor) from the surface of the water to the cathead.
		"I kept her off the wind and sailing free until I had the anchor catted"
[abbreviation]
	- short for computerized axial tomography.
	- computer-assisted (or -aided) testing.
	- clear air turbulence.

$ gtr -f ja -t es 猫
gato
```

### Flags

| Flag     | Default | Description                                    |
| -------- | ------- | ----------------------------------------------- |
| `-f`     | `auto`  | source language code, or `auto` to detect       |
| `-t`     | `en`    | target language code                            |
| `-json`  | `false` | print output as JSON                            |
| `-codes` | `false` | print all available language codes and exit     |

## Install

```sh
go install github.com/jim-ww/gtr@latest
```

### Nix

Try it without installing:

```sh
nix run github:jim-ww/gtr -- <text to translate>
```

Add it to your flake inputs:

```nix
gtr.url = "github:jim-ww/gtr";
```

Then reference `inputs.gtr.packages.${system}.default` wherever you install packages (e.g. `home.packages` or `environment.systemPackages`).

## Support the project

If gtr is useful to you, consider a small donation.

**Monero (XMR)**
```
83YGRqP8uHed6NeegZQeX9ccCxbzoRHHEEi7pTwk4aqdJZEVXXA6NWtetnsEM2v33zFBBt3Rp6DNhU9qhJEGPspU14yN8t7
```

## License

[GPL-3.0](LICENSE). Free to use, study, share, and modify — provided you keep the same freedoms for others.
