monitor PistsVip(SC){

    const max = 30;
    const max_vip = 2; 

    int cantidad_en_pista = 0;
    int cantidad_vip = 0; 


    _condvar cola_esquiadores;
    _condvar cola_vip;

    _proc void ingresarEsquiadorVip(){
        while (cantidad_en_pista == max && cantidad_vip == max_vip) {
            _wait(cola_vip)
        }
        cantidad_en_pista++
        cantidad_vip++ 
    }

    _proc void ingresarEsquiador(){
        while (cantidad_en_pista == max) {
            wait(cola_esquiadores)
        }
        cantidad_en_pista++
    }

    _proc void salirEsquiador(){
        signal(cola_esquiadores)
        cantidad_en_pista--
    }

    _proc void salirEsquiadorVip(){
        signal(cola_vip)
        cantidad_vip--
        signal(cola_esquiadores)
        cantidad_en_pista--
    }

   
}