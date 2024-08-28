Monitor Clinica(SC){

    _condvar cola;   

    bool medico_uno, medico_dos, medico_tres = true; 

   _proc void esperarPaciente(int i){
        wait(paciente);
    }

    _proc void despacharPaciente(int i){
        if i == 1 {
            signal(recepcionista);
            medico_uno = true;
        } else if == 2 {
            signal(recepcionista);
            medico_dos = true;
        } else {
            signal(recepcionista);
            medico_tres = true;
        }
    }

    _proc void recepcionista(){
        wait(recepcionista);
        signal(cola);
    }

    _proc void pacienteEsperar() {
        if !(medico_uno && medico_dos && medico_tres) {
            _wait(cola);
        }
        if medico_uno {
            medico_uno = false;
            signal(paciente);
        } medico_dos {
            medico_dos = false;
            signal(paciente);
        } else {
            medico = false;
            signal(paciente);
        }
    }
}