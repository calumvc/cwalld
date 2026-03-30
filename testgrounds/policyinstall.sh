#!/bin/bash
sudo semodule -i /home/testgrounds/object_types/alpha/alpha.pp
sudo semodule -i /home/testgrounds/object_types/beta/beta.pp
sudo semodule -i /home/testgrounds/object_types/gamma/gamma.pp
sudo semodule -i /home/testgrounds/object_types/sanitised/sanitised.pp

sudo bash /home/testgrounds/subject_types/alpha_rw/alpha_rw.sh
sudo bash /home/testgrounds/subject_types/beta_rw/beta_rw.sh
sudo bash /home/testgrounds/subject_types/gamma_rw/gamma_rw.sh
sudo bash /home/testgrounds/subject_types/alpha_gamma_r/alpha_gamma_r.sh
sudo bash /home/testgrounds/subject_types/beta_gamma_r/beta_gamma_r.sh
sudo bash /home/testgrounds/subject_types/cwalld/cwalld.sh
