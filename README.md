# plcLogger
### Scopo:
Lo scopo di questa app è creare un logger per i dati PLC che sia facile da configurare, scalare e mantenere.

### Idea:
L'idea è di avere una dashboard web per poter visualizzare i dati, modificare le impostazioni e salvare i dati di log sul pc client

### Integrazioni future:
Si potrebbe anche pensare di integrare della API REST per poter rendere accessibili i dati ad altri programmi

### Obbiettivi:
- L'app deve essere super portabile, per questo il linguaggio che verrà usato è GO:
    l'app verrà distribuita come un singolo binario, all'avvio l'app in automatico cercherà il file di configurazione, se presente lo caricherà in memoria, se non è presente verrà creato con dei dati di default.
- L'app deve essere facile da configurare, quindi l'interfaccia web deve essere semplice da usare
- I valori da loggare devono essere impostabili dall'interfaccia web o da config file
- I log saranno di 2 tipi:
    - *Log periodici*: verrà salvato il valore dei dati impostati allo scadere di ogni intervallo di tempo. Il primo intervallo è all'avvio dell'app poi da li parte il conteggio del tempo.
    - *Log on change*: ogni volta che il valore di una tag impostata viene cambiato verrà aggiunta una riga al file di log con timestamp, nome, vallore vecchio e valore nuovo.

### Struttura file:
- il file di congfigurazione è *config.yaml* dove all'interno vengono salvate tutte le impostazioni dell'app
- nella cartella */data* viene salvato il file *last_values.json* con l'ultimo valore dei dati da loggare on change, questo file viene caricato in memoria all'avvio dell'app e poi ogni volta che un valore cambia viene aggiornato.
- nella cartella */log* ci sranno salvati tutti i file di log
    - *onChange.log*: il log dei valori da salvare on change, il file è nella struttura *NDJSON* ogni riga equivale ad un entry json con timestamp, nome, valroe vecchio e valore nuovo.
    - *periodic.log*: il log periodico, anche questo file è nella struttura *NDJSON*, ogni riga equivale ad un log periodico con un timestamp e tutte le tag con il loro valore.

## Cose da fare:
 - [x] Lettura configurazione da file e creazione file default se non esiste
 - [x] Lettura file ultimi valori per store degli utlimi valori salvati per il log sul cambio valore
 - [x] Creazione file di log: periodici e onChange
 - [X] Rotazione file di log se dimensione troppo grande e controllo dimensione archivio per mantenere dimensioni impostate
 - [ ] Lettura valori attivi da PLC {Protocollo S7 con GOS7 o OPCUA}
 - [ ] Creazione interfaccia web essenziale {framework GIN + GO templates + HTMX}
 - [ ] Integrazione gestion archivio file da interfaccia web con possibilità di scaricare i file e cancellarli
 - [ ] Integrazione API REST per recupero dati attivi e storico

 ### Struttura UI
- Una pagina dashboard per vedere:
    - lo stato di connessione del plc
    - i valori show on dashboard
- Una pagina per gestire la configurazione dell'app e la connessione con il plc
- Una pagina per vedere la lista di tag da loggare con i flag in base al log che va fatto e la possibilità di aggiungere nuove tag da loggare
- una pagina per leggere i log più recenti (una per i periodic e una per on change)
- una pagina dove gestire i file di log archiviati
    - poterli leggere
    - poterli scaricare
    - poterli cancellare
- una pagina con lo stato del sistema
    - carico cpu
    - stato ram
    - stato spazio di archiviazione
    - log degli errori

### Optional da implementare nella connessione con PLC
- tentativi di riconnessione nel caso si scolleghi durante la lettura
- scrittura dei dati nel plc
- lettura dei dati a batch dai db
- lettura di dati al di fuori dei db
