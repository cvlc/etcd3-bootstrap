## Create systemd unit
cat << EOT > /etc/systemd/system/ebs-bootstrap-${ebs_volume_name}.service
${ebs_bootstrap_unit} 
EOT

## Enable and start service
systemctl enable ebs-bootstrap-${ebs_volume_name}.service
systemctl start ebs-bootstrap-${ebs_volume_name}.service
