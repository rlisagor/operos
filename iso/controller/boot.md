# Boot sequence for the controller:

1. First boot / static settings
    - bootstrap-settings (first boot only):  apply settings necessary to start
      etcd/Teamster, and those settings which cannot be changed via Teamster
    - ceph-mon-init (first boot only)
    - prepare-addons (first boot only): generate secrets manifests for Kube
      services
2. Teamster initialization
    - etcd
    - operos-cfg-populate (first boot only): pull `/etc/paxautoma/initial-settings`
      into etcd
    - teamster
3. Post-boot settings update
    - start-addons (watch `/var/operos/addons` for kube manifests)
    - apply-settings (run every 10s, update config from Teamster)
