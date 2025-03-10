# eule-booter

I have a couch gaming PC and her name is Eule.
There is no keyboard or mouse connected to Eule, I usually wake her up from a Bluetooth controller.

Eule dual-boots Bazzite and Windows, so I needed an easy way to pick a boot option.

## GRUB

Bazzite uses GRUB, so this solution relies on that.

I added this to `/etc/grub.d/40_custom`:

```shell
insmod net
insmod efinet
insmod tftp

net_bootp

source (tftp,eule-booter.lan.al)
```

This makes GRUB initialise network, receive configuration over DHCP, and then attempt sourcing a config from a remote TFTP server. `eule-booter.lan.al` is a DNS name of an LXC running software from this repo.

## TFTP

This software implements a very simple TFTP server. It will always serve this, where `%s` is the **next boot option**:

```shell
set default="%s"
```

By default, the next boot option is **0**. It can be updated over HTTP.

## HTTP

This software provides a basic HTTP API available on port 80.

Get current next boot option:

```shell
curl http://eule-booter.lan.al
```

Set next boot option (e.g. `uefi-firmware`):
```shell
curl -d 'uefi-firmware' http://eule-booter.lan.al
```

eule-booter keeps the next boot option in RAM, so it will return back to **0** upon restart.

## Boot options

Next boot option can be anything your GRUB setup can handle in the `default` variable.

It could be an ID of a menuentry (recommended):
- `ostree-0-9b67f3ea-e3d2-4450-8c8e-cab650855c84`
- `osprober-efi-BBD8-F97D`
- `uefi-firmware`

It could be a full name of a menuentry (meh):
- `Bazzite 41 (FROM Fedora Kinoite) (ostree:0)`
- `Windows Boot Manager (on /dev/nvme1n1p1)`
- `UEFI Firmware Settings`

It could also be a zero-based index of a menuentry (guh).

Check your GRUB config to see available menu entries and how to reference them.

There might be other ways to reference menu entries I'm not aware of.