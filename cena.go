package main

    import (
        "hash/fnv"
        "log"
        "math/rand"
        "os"
        "time"
    )

    // el numero de filosofos es el tamaño de la lista.
    var ph = []string{"Cristiano Ronaldo", "Mesii", "Silverter Staloone", "Max", "Russell"}

    const hambriento = 3                // numero de tiempos que cada filosofos come.
    const pensar = time.Second / 100 // la hora de pensar
    const comer = time.Second / 100   // la hhora de comer

    var fmt = log.New(os.Stdout, "", 0) // para cada hilo seguro de  salida.

    var hecho = make(chan bool)

    // se usan canales para implementar la solucion.
    // enviar sobre los canales los "tenedores."
    type fork byte

    // un tenedoren el programa modela un tenedor fisico en la c¡simulacion.
    // un canal separado representa cada espacio para el tenedor entre dos filosofos.
    // se tienne acceso a cada tenedor. los canales son reducidos a capacidad = 1.

    // Gorutina paara las acciones de los filosofos. una instancia se ejecuta
    // para cada filosofo.  Las instancias se ejeuctan concurrentemente.
    func filosofoo(NombreFilosofo string,
        ManoDominante, otraMano chan fork, hecho chan bool) {
        fmt.Println(NombreFilosofo, "sentado")
        // cada Gorutina del filosofo un número aleatorio,
        // con un hash de los nombres de los filosofos.
        h := fnv.New64a()
        h.Write([]byte(NombreFilosofo))
        rg := rand.New(rand.NewSource(int64(h.Sum64())))
        // funcion util como ejecucion sleep para
        rSleep := func(t time.Duration) {
            time.Sleep(t + time.Duration(rg.Int63n(int64(t))))
        }
        for h := hambriento; h > 0; h-- {
            fmt.Println(NombreFilosofo, "hambruna")
            <-ManoDominante // levanta los tenedores
            <-otraMano
            fmt.Println(NombreFilosofo, "comiendo")
            rSleep(comer)
            ManoDominante <- 'f' // deja los tenedores
            otraMano <- 'f'
            fmt.Println(NombreFilosofo, "filosofando")
            rSleep(pensar)
        }
        fmt.Println(NombreFilosofo, "satisfecho")
        hecho <- true
        fmt.Println(NombreFilosofo, "se va de la mesa")
    }

    func main() {
        fmt.Println("la mesa esta vacía")
        // creaar los canales de los tenedores y se inicia las Gorutinas de los filosofos.
        // cada Gorutina con los canales
        lugar := make(chan fork, 1)
        lugar <- 'f' // el byte en el canal representa un tenedor en la mesa.
        lugarIzquierda := lugar
        for i := 1; i < len(ph); i++ {
            lugarDerecha := make(chan fork, 1)
            lugarDerecha <- 'f'
            go filosofoo(ph[i], lugarIzquierda, lugarDerecha, hecho)
            lugarIzquierda = lugarDerecha
        }
        // hacer un unico filosofo su mano izquierda by reversing fork place
        // supplied to filosofoo's dominant hand.
        // esto hace precedencia aciclico,  preventing deadlock.
        go filosofoo(ph[0], lugar, lugarIzquierda, hecho)
        // ellos estan ocupados comiendo ahora.
        for range ph {
            <-hecho // esperar a los filosofos para finalizar.
        }
        fmt.Println("la mesa está vacía")
    }
