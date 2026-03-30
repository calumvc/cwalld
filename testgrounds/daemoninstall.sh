#!/bin/bash

gcc -o /home/testgrounds/subjects/alpha_civil/alpha_civild /home/testgrounds/subjects/alpha_civil/alpha_civild.c
sudo cp /home/testgrounds/subjects/alpha_civil/alpha_civild /usr/local/bin
sudo cp /home/testgrounds/subjects/alpha_civil/alpha_civild.service /usr/lib/systemd/system

gcc -o /home/testgrounds/subjects/beta_civil/beta_civild /home/testgrounds/subjects/beta_civil/beta_civild.c
sudo cp /home/testgrounds/subjects/beta_civil/beta_civild /usr/local/bin
sudo cp /home/testgrounds/subjects/beta_civil/beta_civild.service /usr/lib/systemd/system

gcc -o /home/testgrounds/subjects/gamma_civil/gamma_civild /home/testgrounds/subjects/gamma_civil/gamma_civild.c
sudo cp /home/testgrounds/subjects/gamma_civil/gamma_civild /usr/local/bin
sudo cp /home/testgrounds/subjects/gamma_civil/gamma_civild.service /usr/lib/systemd/system

gcc -o /home/testgrounds/subjects/alpha_curious/alpha_curiousd /home/testgrounds/subjects/alpha_curious/alpha_curiousd.c
sudo cp /home/testgrounds/subjects/alpha_curious/alpha_curiousd /usr/local/bin
sudo cp /home/testgrounds/subjects/alpha_curious/alpha_curiousd.service /usr/lib/systemd/system

gcc -o /home/testgrounds/subjects/beta_evil/beta_evild /home/testgrounds/subjects/beta_evil/beta_evild.c
sudo cp /home/testgrounds/subjects/beta_evil/beta_evild /usr/local/bin
sudo cp /home/testgrounds/subjects/beta_evil/beta_evild.service /usr/lib/systemd/system

gcc -o /home/testgrounds/subjects/gamma_confused/gamma_confusedd /home/testgrounds/subjects/gamma_confused/gamma_confusedd.c
sudo cp /home/testgrounds/subjects/gamma_confused/gamma_confusedd /usr/local/bin
sudo cp /home/testgrounds/subjects/gamma_confused/gamma_confusedd.service /usr/lib/systemd/system

gcc -o /home/testgrounds/subjects/unconfined_daemon/unconfined_daemond /home/testgrounds/subjects/unconfined_daemon/unconfined_daemond.c
sudo cp /home/testgrounds/subjects/unconfined_daemon/unconfined_daemond /usr/local/bin
sudo cp /home/testgrounds/subjects/unconfined_daemon/unconfined_daemond.service /usr/lib/systemd/system

sudo systemctl daemon-reload
