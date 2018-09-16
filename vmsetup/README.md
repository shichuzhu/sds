# Setting up the vms

## To setup the vm for the first time

1. Create ppk pair with `ssh-keygen`
2. Tune user (netid) and group number info in `deploy` file.
3. Deploy the public keys to vms. `./deploy CopyKey`
   - Need to "yes" and "passwd"
4. For each vm, ssh and run the setup script. `./deploy For 2 Each "bash -s" '< setup.sh'`
5. To update source code on VMs upon testing new code, use <br />
`vmsetup/deploy Copy src @go/src/fa18cs425mp/`
<br />Note this does NOT copy the Makefile.

