#!/bin/bash

sudo chcon -t alpha_t /home/testgrounds/objects/alpha_logs
sudo chcon -t beta_t /home/testgrounds/objects/beta_plans
sudo chcon -t gamma_t /home/testgrounds/objects/gamma_reports
sudo chcon -t sanitised_t /home/testgrounds/objects/sanitised

sudo chcon -t bin_t /usr/local/bin/*
sudo chcon -u system_u /usr/local/bin/*

sudo semanage fcontext -a -t cwalld_exec_t /usr/local/sbin/cwalld-enforce
sudo restorecon -v /usr/local/sbin/cwalld-enforce
