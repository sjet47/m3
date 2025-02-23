# m<sup>3</sup>: A Minecraft Mod Manager

Download/Update/Manage your Minecraft mods in different version all together

## Install

```shell
go install github.com/sjet47/m3
```

## Usage

### Acquire CurseForge API Key

To use `m3`, you must have an available CurseForge API key. And you need to set the `CURSE_FORGE_APIKEY` environment variable to your CurseForge API key.
To get the API key, you need a [CurseForge For Studios](https://console.curseforge.com/#/) account and generate an API key from [there](https://console.curseforge.com/#/api-keys).

Then you can export this environment variable in your shell configuration file.

For Bash:

```
echo export CURSE_FORGE_APIKEY='your_api_key' >> ~/.bashrc
```

For ZSH:

```
echo export CURSE_FORGE_APIKEY='your_api_key' >> ~/.zshrc
```

Or pass the environment variable directly each time running `m3`(which is more annoying but more secure):

```shell
CURSE_FORGE_APIKEY='your_api_key' m3 <subcommand>...
```

### Init

First `cd` to the `mods` directory that contains the mods you want to manage, e.g. `.minecraft/mods` or `.minecraft/versions/<version>/mods`


```shell
m3 init <minecraft_version>
```

> It's recommended to use a separate directory for each version of Minecraft, which can be done by setting the **Game Directory** in the installation profile of the Minecraft Launcher, so you can isolate the `mods` directory from different versions of Minecraft.

### Add Mods

Add mods into `m3`'s index.

```shell
m3 add <mod_id>... [-l mod_loader] [-o]
```

- `-l --option`: Specify the mod loader, e.g. `Forge`, `Fabric`, `NeoForge`, etc. The name is not case sensitive. Default mod loader is `Forge`.
- `-o --option`: Whether to download optional dependencies. Only necessary dependencies will be downloaded by default.

You can also use file to provide mod ids with `-f` flag when there are too many mods to add at once.

```shell
m3 add -f <file> [-l mod_loader] [-o]
```

The file should be a csv format file without header line and mod id is in the first column. All other columns will be ignored. A example file is like this:

```csv
238222,JET
240630,JER
32274,Journey Map
248787,AppleSkin
```

### Update Mods

Update mods in `m3`'s index.

```shell
m3 update
```

This command will only update the mods with the same mod loader and whether to download optional dependencies as `m3 add`. To change that, you can rerun `m3 add` to update the mod index.

### List Mods In Index

Wanna know what mods are in the index? Just run:

```shell
m3 ls
```

### Remove Mods From Index

If you decide to remove some mods from the index, you can run:

```shell
m3 remove <mod_id>...
```

The mod files will also be removed from the `mods` directory.
