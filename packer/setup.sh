#!/bin/bash
sudo mv /home/arch/flipbot.service /etc/systemd/system/
sudo mv /home/arch/flipbotreload.service /etc/systemd/system/
mkdir /home/arch/flipbot
sudo systemctl enable flipbot.service
sudo pacman -Syu --noconfirm yay bash-completion ncdu nano git wget mosh curl iperf3
yay -S --noconfirm translate-shell googler
echo "UUID=83282862-5fc8-4e9a-be23-03840738cb2d       /home/arch/flipbot      ext4    rw,relatime     0 2" | sudo tee -a /etc/fstab
echo -e 'en_US.UTF-8 UTF-8\nlv_LV.UTF-8 UTF-8' | sudo tee -a /etc/locale.gen
sudo locale-gen
sudo pacman -Scc --noconfirm