package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Estructuras del sistema

type Cliente struct {
	ID       int
	Nombre   string
	Telefono string
	Email    string
	Vehiculo *Vehiculo
}

type Vehiculo struct {
	Matricula    string
	Marca        string
	Modelo       string
	FechaEntrada string
	FechaSalida  string
	Incidencia   *Incidencia
	EnTaller     bool
	NumeroPlaza  int
}

type Mecanico struct {
	ID           int
	Nombre       string
	Especialidad string
	AniosExp     int
	Activo       bool
	Incidencias  []*Incidencia
}

type Incidencia struct {
	ID          int
	Mecanicos   []*Mecanico
	Tipo        string
	Prioridad   string
	Descripcion string
	Estado      string
}

type Taller struct {
	Mecanicos         []*Mecanico
	PlazasPorMecanico int
	PlazasOcupadas    map[int]bool
	TotalPlazas       int
}

// Variables globales
var (
	clientes    []*Cliente
	vehiculos   []*Vehiculo
	incidencias []*Incidencia
	mecanicos   []*Mecanico
	taller      Taller

	contadorCliente    = 1
	contadorIncidencia = 1
	contadorMecanico   = 1
)

// Funciones auxiliares

func limpiarPantalla() {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		cmd = exec.Command("clear")
	} else {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func pausar() {
	fmt.Println("\nPresione ENTER para continuar...")
	fmt.Scanln()
}

func obtenerFechaActual() string {
	return time.Now().Format("02/01/2006")
}

func calcularTotalPlazas() int {
	total := 0
	for _, m := range mecanicos {
		if m.Activo {
			total += taller.PlazasPorMecanico
		}
	}
	return total
}

func contarPlazasOcupadas() int {
	contador := 0
	for _, ocupada := range taller.PlazasOcupadas {
		if ocupada {
			contador++
		}
	}
	return contador
}

// 1. Funciones CRUD - Gestion de Clientes

func crearCliente() {
	limpiarPantalla()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== CREAR CLIENTE ===")

	var telefono, email string

	fmt.Print("Nombre del cliente: ")
	nombre, _ := reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	if nombre == "" {
		fmt.Println("Error: El nombre no puede estar vacío")
		pausar()
		return
	}

	fmt.Print("Teléfono: ")
	fmt.Scanln(&telefono)

	fmt.Print("Email: ")
	fmt.Scanln(&email)

	cliente := &Cliente{
		ID:       contadorCliente,
		Nombre:   nombre,
		Telefono: telefono,
		Email:    email,
		Vehiculo: nil,
	}

	clientes = append(clientes, cliente)
	contadorCliente++

	fmt.Println("\nCliente creado exitosamente con ID:", cliente.ID)
	pausar()
}

func visualizarClientes() {
	limpiarPantalla()
	fmt.Println("=== LISTA DE CLIENTES ===")

	if len(clientes) == 0 {
		fmt.Println("No hay clientes registrados")
		pausar()
		return
	}

	for _, c := range clientes {
		fmt.Printf("\nID: %d\n", c.ID)
		fmt.Printf("Nombre: %s\n", c.Nombre)
		fmt.Printf("Teléfono: %s\n", c.Telefono)
		fmt.Printf("Email: %s\n", c.Email)
		if c.Vehiculo != nil {
			fmt.Printf("Vehículo asociado: %s (Matrícula: %s)\n",
				c.Vehiculo.Marca+" "+c.Vehiculo.Modelo, c.Vehiculo.Matricula)
		} else {
			fmt.Println("Sin vehículo asociado")
		}
		fmt.Println("---")
	}

	pausar()
}

func modificarCliente() {
	limpiarPantalla()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== MODIFICAR CLIENTE ===")

	var id int
	fmt.Print("ID del cliente a modificar: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	var cliente *Cliente
	for _, c := range clientes {
		if c.ID == id {
			cliente = c
			break
		}
	}

	if cliente == nil {
		fmt.Println("Error: Cliente no encontrado")
		pausar()
		return
	}

	fmt.Printf("\nCliente actual: %s\n", cliente.Nombre)
	fmt.Print("Nuevo nombre (dejar vacío para no cambiar): ")
	nombre, _ := reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)
	if nombre != "" {
		cliente.Nombre = nombre
	}

	fmt.Print("Nuevo teléfono (dejar vacío para no cambiar): ")
	var telefono string
	fmt.Scanln(&telefono)
	if telefono != "" {
		cliente.Telefono = telefono
	}

	fmt.Print("Nuevo email (dejar vacío para no cambiar): ")
	var email string
	fmt.Scanln(&email)
	if email != "" {
		cliente.Email = email
	}

	fmt.Println("\nCliente modificado exitosamente")
	pausar()
}

func eliminarCliente() {
	limpiarPantalla()
	fmt.Println("=== ELIMINAR CLIENTE ===")

	var id int
	fmt.Print("ID del cliente a eliminar: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	for i, c := range clientes {
		// Si se encuentra el cliente
		if c.ID == id {
			// Si tiene vehículo, liberarlo del taller
			if c.Vehiculo != nil && c.Vehiculo.EnTaller {
				taller.PlazasOcupadas[c.Vehiculo.NumeroPlaza] = false
			}
			// Eliminar el cliente
			clientes = append(clientes[:i], clientes[i+1:]...)
			fmt.Println("Cliente eliminado exitosamente")
			pausar()
			return
		}
	}

	fmt.Println("Error: Cliente no encontrado")
	pausar()
}

// 2. Funciones CRUD - Gestion de Vehículos

func crearVehiculo() {
	limpiarPantalla()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== CREAR VEHÍCULO ===")

	var idCliente int
	fmt.Print("ID del cliente propietario: ")
	fmt.Scanf("%d", &idCliente)
	fmt.Scanln()

	var cliente *Cliente
	for _, c := range clientes {
		if c.ID == idCliente {
			cliente = c
			break
		}
	}

	if cliente == nil {
		fmt.Println("Error: Cliente no encontrado")
		pausar()
		return
	}

	if cliente.Vehiculo != nil {
		fmt.Println("Error: El cliente ya tiene un vehículo asociado")
		pausar()
		return
	}

	fmt.Print("Matricula: ")
	matricula, _ := reader.ReadString('\n')
	matricula = strings.TrimSpace(matricula)

	if matricula == "" {
		fmt.Println("Error: La matrícula no puede estar vacía")
		pausar()
		return
	}

	// Verificar que la matrícula no exista
	for _, v := range vehiculos {
		if v.Matricula == matricula {
			fmt.Println("Error: Ya existe un vehículo con esa matrícula")
			pausar()
			return
		}
	}

	fmt.Print("Marca: ")
	marca, _ := reader.ReadString('\n')
	marca = strings.TrimSpace(marca)

	fmt.Print("Modelo: ")
	modelo, _ := reader.ReadString('\n')
	modelo = strings.TrimSpace(modelo)

	vehiculo := &Vehiculo{
		Matricula:    matricula,
		Marca:        marca,
		Modelo:       modelo,
		FechaEntrada: obtenerFechaActual(),
		FechaSalida:  "",
		Incidencia:   nil,
		EnTaller:     false,
		NumeroPlaza:  -1,
	}

	vehiculos = append(vehiculos, vehiculo)
	cliente.Vehiculo = vehiculo

	fmt.Println("\nVehículo creado exitosamente y asociado al cliente")
	pausar()
}

func visualizarVehiculos() {
	limpiarPantalla()
	fmt.Println("=== LISTA DE VEHÍCULOS ===")

	if len(vehiculos) == 0 {
		fmt.Println("No hay vehículos registrados")
		pausar()
		return
	}

	for _, v := range vehiculos {
		fmt.Printf("\nMatrícula: %s\n", v.Matricula)
		fmt.Printf("Marca: %s\n", v.Marca)
		fmt.Printf("Modelo: %s\n", v.Modelo)
		fmt.Printf("Fecha entrada: %s\n", v.FechaEntrada)
		if v.FechaSalida != "" {
			fmt.Printf("Fecha salida estimada: %s\n", v.FechaSalida)
		}
		if v.EnTaller {
			fmt.Printf("En taller: Sí (Plaza %d)\n", v.NumeroPlaza)
		} else {
			fmt.Println("En taller: No")
		}
		if v.Incidencia != nil {
			fmt.Printf("Incidencia: ID %d - %s (%s)\n",
				v.Incidencia.ID, v.Incidencia.Tipo, v.Incidencia.Estado)
		}
		fmt.Println("---")
	}

	pausar()
}

func modificarVehiculo() {
	limpiarPantalla()
	fmt.Println("=== MODIFICAR VEHÍCULO ===")

	var matricula string
	fmt.Print("Matrícula del vehículo a modificar: ")
	fmt.Scanln(&matricula)

	var vehiculo *Vehiculo
	for _, v := range vehiculos {
		if v.Matricula == matricula {
			vehiculo = v
			break
		}
	}

	if vehiculo == nil {
		fmt.Println("Error: Vehículo no encontrado")
		pausar()
		return
	}

	fmt.Printf("\nVehículo actual: %s %s\n", vehiculo.Marca, vehiculo.Modelo)

	fmt.Print("Nueva marca (dejar vacío para no cambiar): ")
	var marca string
	fmt.Scanln(&marca)
	if marca != "" {
		vehiculo.Marca = marca
	}

	fmt.Print("Nuevo modelo (dejar vacío para no cambiar): ")
	var modelo string
	fmt.Scanln(&modelo)
	if modelo != "" {
		vehiculo.Modelo = modelo
	}

	fmt.Print("Nueva fecha salida estimada (DD/MM/AAAA, vacío para no cambiar): ")
	var fecha string
	fmt.Scanln(&fecha)
	if fecha != "" {
		vehiculo.FechaSalida = fecha
	}

	fmt.Println("\nVehículo modificado exitosamente")
	pausar()
}

func eliminarVehiculo() {
	limpiarPantalla()
	fmt.Println("=== ELIMINAR VEHÍCULO ===")

	var matricula string
	fmt.Print("Matrícula del vehículo a eliminar: ")
	fmt.Scanln(&matricula)

	for i, v := range vehiculos {
		if v.Matricula == matricula {
			// Liberar plaza si está en taller
			if v.EnTaller {
				taller.PlazasOcupadas[v.NumeroPlaza] = false
			}

			// Desvincular del cliente
			for _, c := range clientes {
				if c.Vehiculo == v {
					c.Vehiculo = nil
					break
				}
			}

			vehiculos = append(vehiculos[:i], vehiculos[i+1:]...)
			fmt.Println("Vehículo eliminado exitosamente")
			pausar()
			return
		}
	}

	fmt.Println("Error: Vehículo no encontrado")
	pausar()
}

// 3. Funciones CRUD - Gestion de Incidencias

func crearIncidencia() {
	limpiarPantalla()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== CREAR INCIDENCIA ===")

	fmt.Print("Matricula: ")
	matricula, _ := reader.ReadString('\n')
	matricula = strings.TrimSpace(matricula)

	var vehiculo *Vehiculo
	for _, v := range vehiculos {
		if v.Matricula == matricula {
			vehiculo = v
			break
		}
	}

	if vehiculo == nil {
		fmt.Println("Error: Vehículo no encontrado")
		pausar()
		return
	}

	if vehiculo.Incidencia != nil {
		fmt.Println("Error: El vehículo ya tiene una incidencia asignada")
		pausar()
		return
	}

	var tipo, prioridad string

	fmt.Println("\nTipo de incidencia:")
	fmt.Println("1. Mecánica")
	fmt.Println("2. Eléctrica")
	fmt.Println("3. Carrocería")
	var opcion int
	fmt.Print("Opción: ")
	fmt.Scanf("%d", &opcion)
	fmt.Scanln()

	switch opcion {
	case 1:
		tipo = "mecánica"
	case 2:
		tipo = "eléctrica"
	case 3:
		tipo = "carrocería"
	default:
		fmt.Println("Opción inválida")
		pausar()
		return
	}

	fmt.Println("\nPrioridad:")
	fmt.Println("1. Baja")
	fmt.Println("2. Media")
	fmt.Println("3. Alta")
	fmt.Print("Opción: ")
	fmt.Scanf("%d", &opcion)
	fmt.Scanln()

	switch opcion {
	case 1:
		prioridad = "baja"
	case 2:
		prioridad = "media"
	case 3:
		prioridad = "alta"
	default:
		fmt.Println("Opción inválida")
		pausar()
		return
	}

	fmt.Print("Descripción: ")
	descripcion, _ := reader.ReadString('\n')
	descripcion = strings.TrimSpace(descripcion)

	incidencia := &Incidencia{
		ID:          contadorIncidencia,
		Mecanicos:   []*Mecanico{},
		Tipo:        tipo,
		Prioridad:   prioridad,
		Descripcion: descripcion,
		Estado:      "abierta",
	}

	incidencias = append(incidencias, incidencia)
	vehiculo.Incidencia = incidencia
	contadorIncidencia++

	fmt.Println("\nIncidencia creada exitosamente con ID:", incidencia.ID)
	pausar()
}

func visualizarIncidencias() {
	limpiarPantalla()
	fmt.Println("=== LISTA DE INCIDENCIAS ===")

	if len(incidencias) == 0 {
		fmt.Println("No hay incidencias registradas")
		pausar()
		return
	}

	for _, inc := range incidencias {
		fmt.Printf("\nID: %d\n", inc.ID)
		fmt.Printf("Tipo: %s\n", inc.Tipo)
		fmt.Printf("Prioridad: %s\n", inc.Prioridad)
		fmt.Printf("Estado: %s\n", inc.Estado)
		fmt.Printf("Descripción: %s\n", inc.Descripcion)

		if len(inc.Mecanicos) > 0 {
			fmt.Print("Mecánicos asignados: ")
			for i, m := range inc.Mecanicos {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s", m.Nombre)
			}
			fmt.Println()
		} else {
			fmt.Println("Sin mecánicos asignados")
		}
		fmt.Println("---")
	}

	pausar()
}

func modificarIncidencia() {
	limpiarPantalla()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== MODIFICAR INCIDENCIA ===")

	var id int
	fmt.Print("ID de la incidencia a modificar: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	var incidencia *Incidencia
	for _, inc := range incidencias {
		if inc.ID == id {
			incidencia = inc
			break
		}
	}

	if incidencia == nil {
		fmt.Println("Error: Incidencia no encontrada")
		pausar()
		return
	}

	fmt.Printf("\nIncidencia actual: %s (%s)\n", incidencia.Tipo, incidencia.Estado)

	fmt.Print("Nueva descripción (dejar vacío para no cambiar): ")
	descripcion, _ := reader.ReadString('\n')
	descripcion = strings.TrimSpace(descripcion)
	if descripcion != "" {
		incidencia.Descripcion = descripcion
	}

	fmt.Println("\nCambiar prioridad? (S/N): ")
	var cambiar string
	fmt.Scanln(&cambiar)
	if strings.ToUpper(cambiar) == "S" {
		fmt.Println("1. Baja")
		fmt.Println("2. Media")
		fmt.Println("3. Alta")
		var opcion int
		fmt.Print("Opción: ")
		fmt.Scanf("%d", &opcion)
		fmt.Scanln()

		switch opcion {
		case 1:
			incidencia.Prioridad = "baja"
		case 2:
			incidencia.Prioridad = "media"
		case 3:
			incidencia.Prioridad = "alta"
		}
	}

	fmt.Println("\nIncidencia modificada exitosamente")
	pausar()
}

func eliminarIncidencia() {
	limpiarPantalla()
	fmt.Println("=== ELIMINAR INCIDENCIA ===")

	var id int
	fmt.Print("ID de la incidencia a eliminar: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	for i, inc := range incidencias {
		if inc.ID == id {
			// Desvincular de vehículo
			for _, v := range vehiculos {
				if v.Incidencia == inc {
					v.Incidencia = nil
					break
				}
			}

			// Desvincular de mecánicos
			for _, m := range inc.Mecanicos {
				for j, incAsig := range m.Incidencias {
					if incAsig == inc {
						m.Incidencias = append(m.Incidencias[:j], m.Incidencias[j+1:]...)
						break
					}
				}
			}

			incidencias = append(incidencias[:i], incidencias[i+1:]...)
			fmt.Println("Incidencia eliminada exitosamente")
			pausar()
			return
		}
	}

	fmt.Println("Error: Incidencia no encontrada")
	pausar()
}

// 4. Funciones CRUD - Gestion de Mecánicos

func crearMecanico() {
	limpiarPantalla()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== CREAR MECÁNICO ===")

	var especialidad string
	var anios int

	fmt.Print("Nombre del mecánico: ")
	nombre, _ := reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	if nombre == "" {
		fmt.Println("Error: El nombre no puede estar vacío")
		pausar()
		return
	}

	fmt.Println("\nEspecialidad:")
	fmt.Println("1. Mecánica")
	fmt.Println("2. Eléctrica")
	fmt.Println("3. Carrocería")
	var opcion int
	fmt.Print("Opción: ")
	fmt.Scanf("%d", &opcion)
	fmt.Scanln()

	switch opcion {
	case 1:
		especialidad = "mecánica"
	case 2:
		especialidad = "eléctrica"
	case 3:
		especialidad = "carrocería"
	default:
		fmt.Println("Opción inválida")
		pausar()
		return
	}

	fmt.Print("Años de experiencia: ")
	fmt.Scanf("%d", &anios)
	fmt.Scanln()

	mecanico := &Mecanico{
		ID:           contadorMecanico,
		Nombre:       nombre,
		Especialidad: especialidad,
		AniosExp:     anios,
		Activo:       true,
		Incidencias:  []*Incidencia{},
	}

	mecanicos = append(mecanicos, mecanico)
	taller.Mecanicos = append(taller.Mecanicos, mecanico)
	taller.TotalPlazas = calcularTotalPlazas()
	contadorMecanico++

	fmt.Println("\nMecánico creado exitosamente con ID:", mecanico.ID)
	pausar()
}

func visualizarMecanicos() {
	limpiarPantalla()
	fmt.Println("=== LISTA DE MECÁNICOS ===")

	if len(mecanicos) == 0 {
		fmt.Println("No hay mecánicos registrados")
		pausar()
		return
	}

	for _, m := range mecanicos {
		fmt.Printf("\nID: %d\n", m.ID)
		fmt.Printf("Nombre: %s\n", m.Nombre)
		fmt.Printf("Especialidad: %s\n", m.Especialidad)
		fmt.Printf("Años de experiencia: %d\n", m.AniosExp)
		if m.Activo {
			fmt.Println("Estado: Activo")
		} else {
			fmt.Println("Estado: De baja")
		}
		fmt.Printf("Incidencias asignadas: %d\n", len(m.Incidencias))
		fmt.Println("---")
	}

	pausar()
}

func modificarMecanico() {
	limpiarPantalla()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== MODIFICAR MECÁNICO ===")

	var id int
	fmt.Print("ID del mecánico a modificar: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	var mecanico *Mecanico
	for _, m := range mecanicos {
		if m.ID == id {
			mecanico = m
			break
		}
	}

	if mecanico == nil {
		fmt.Println("Error: Mecánico no encontrado")
		pausar()
		return
	}

	fmt.Printf("\nMecánico actual: %s\n", mecanico.Nombre)

	fmt.Print("Nuevo nombre (dejar vacío para no cambiar): ")
	nombre, _ := reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)
	if nombre != "" {
		mecanico.Nombre = nombre
	}

	fmt.Print("Nuevos años de experiencia (0 para no cambiar): ")
	var anios int
	fmt.Scanf("%d", &anios)
	fmt.Scanln()
	if anios > 0 {
		mecanico.AniosExp = anios
	}

	fmt.Println("\nMecánico modificado exitosamente")
	pausar()
}

func eliminarMecanico() {
	limpiarPantalla()
	fmt.Println("=== ELIMINAR MECÁNICO ===")

	var id int
	fmt.Print("ID del mecánico a eliminar: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	for i, m := range mecanicos {
		if m.ID == id {
			if len(m.Incidencias) > 0 {
				fmt.Println("Error: No se puede eliminar un mecánico con incidencias asignadas")
				pausar()
				return
			}

			mecanicos = append(mecanicos[:i], mecanicos[i+1:]...)

			// Actualizar taller
			for j, mecTaller := range taller.Mecanicos {
				if mecTaller == m {
					taller.Mecanicos = append(taller.Mecanicos[:j], taller.Mecanicos[j+1:]...)
					break
				}
			}

			taller.TotalPlazas = calcularTotalPlazas()

			fmt.Println("Mecánico eliminado exitosamente")
			pausar()
			return
		}
	}

	fmt.Println("Error: Mecánico no encontrado")
	pausar()
}

// Funciones operativas del taller

func asignarVehiculoATaller() {
	limpiarPantalla()
	fmt.Println("=== ASIGNAR VEHÍCULO A TALLER ===")

	var matricula string
	fmt.Print("Matrícula del vehículo: ")
	fmt.Scanln(&matricula)

	var vehiculo *Vehiculo
	for _, v := range vehiculos {
		if v.Matricula == matricula {
			vehiculo = v
			break
		}
	}

	if vehiculo == nil {
		fmt.Println("Error: Vehículo no encontrado")
		pausar()
		return
	}

	if vehiculo.EnTaller {
		fmt.Println("Error: El vehículo ya está en el taller")
		pausar()
		return
	}

	// Verificar plazas disponibles
	plazasOcupadas := contarPlazasOcupadas()
	totalPlazas := calcularTotalPlazas()

	if plazasOcupadas >= totalPlazas {
		fmt.Println("Error: No hay plazas disponibles en el taller")
		pausar()
		return
	}

	// Buscar primera plaza libre
	plazaAsignada := -1
	for i := 1; i <= totalPlazas; i++ {
		if !taller.PlazasOcupadas[i] {
			plazaAsignada = i
			break
		}
	}

	vehiculo.EnTaller = true
	vehiculo.NumeroPlaza = plazaAsignada
	taller.PlazasOcupadas[plazaAsignada] = true

	fmt.Printf("\nVehículo asignado exitosamente a la plaza %d\n", plazaAsignada)
	pausar()
}

func visualizarEstadoTaller() {
	limpiarPantalla()
	fmt.Println("=== ESTADO DEL TALLER ===")

	totalPlazas := calcularTotalPlazas()
	plazasOcupadas := contarPlazasOcupadas()
	plazasLibres := totalPlazas - plazasOcupadas

	fmt.Printf("\nTotal de plazas: %d\n", totalPlazas)
	fmt.Printf("Plazas ocupadas: %d\n", plazasOcupadas)
	fmt.Printf("Plazas libres: %d\n", plazasLibres)

	fmt.Println("\n--- Detalle de plazas ocupadas ---")
	for i := 1; i <= totalPlazas; i++ {
		if taller.PlazasOcupadas[i] {
			for _, v := range vehiculos {
				if v.NumeroPlaza == i {
					fmt.Printf("Plaza %d: %s %s (Matrícula: %s)\n",
						i, v.Marca, v.Modelo, v.Matricula)
					break
				}
			}
		}
	}

	fmt.Println("\n--- Mecánicos activos ---")
	for _, m := range mecanicos {
		if m.Activo {
			fmt.Printf("%s (%s) - %d incidencias asignadas\n",
				m.Nombre, m.Especialidad, len(m.Incidencias))
		}
	}

	pausar()
}

func darAltaBajaMecanico() {
	limpiarPantalla()
	fmt.Println("=== DAR ALTA/BAJA A MECÁNICO ===")

	var id int
	fmt.Print("ID del mecánico: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	var mecanico *Mecanico
	for _, m := range mecanicos {
		if m.ID == id {
			mecanico = m
			break
		}
	}

	if mecanico == nil {
		fmt.Println("Error: Mecánico no encontrado")
		pausar()
		return
	}

	if mecanico.Activo {
		if len(mecanico.Incidencias) > 0 {
			fmt.Println("Error: No se puede dar de baja a un mecánico con incidencias asignadas")
			pausar()
			return
		}
		mecanico.Activo = false
		fmt.Println("Mecánico dado de baja exitosamente")
	} else {
		mecanico.Activo = true
		fmt.Println("Mecánico dado de alta exitosamente")
	}

	taller.TotalPlazas = calcularTotalPlazas()
	pausar()
}

func cambiarEstadoIncidencia() {
	limpiarPantalla()
	fmt.Println("=== CAMBIAR ESTADO DE INCIDENCIA ===")

	var id int
	fmt.Print("ID de la incidencia: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	var incidencia *Incidencia
	for _, inc := range incidencias {
		if inc.ID == id {
			incidencia = inc
			break
		}
	}

	if incidencia == nil {
		fmt.Println("Error: Incidencia no encontrada")
		pausar()
		return
	}

	fmt.Printf("\nEstado actual: %s\n", incidencia.Estado)
	fmt.Println("\nNuevo estado:")
	fmt.Println("1. Abierta")
	fmt.Println("2. En proceso")
	fmt.Println("3. Cerrada")

	var opcion int
	fmt.Print("Opción: ")
	fmt.Scanf("%d", &opcion)
	fmt.Scanln()

	switch opcion {
	case 1:
		incidencia.Estado = "abierta"
	case 2:
		incidencia.Estado = "en proceso"
	case 3:
		incidencia.Estado = "cerrada"
		// Si se cierra, liberar el vehículo del taller
		for _, v := range vehiculos {
			if v.Incidencia == incidencia && v.EnTaller {
				v.EnTaller = false
				taller.PlazasOcupadas[v.NumeroPlaza] = false
				v.NumeroPlaza = -1
				fmt.Println("El vehículo ha sido liberado del taller")
			}
		}
	default:
		fmt.Println("Opción inválida")
		pausar()
		return
	}

	fmt.Println("\nEstado de incidencia cambiado exitosamente")
	pausar()
}

func asignarMecanicoAIncidencia() {
	limpiarPantalla()
	fmt.Println("=== ASIGNAR MECÁNICO A INCIDENCIA ===")

	var idIncidencia int
	fmt.Print("ID de la incidencia: ")
	fmt.Scanf("%d", &idIncidencia)
	fmt.Scanln()

	var incidencia *Incidencia
	for _, inc := range incidencias {
		if inc.ID == idIncidencia {
			incidencia = inc
			break
		}
	}

	if incidencia == nil {
		fmt.Println("Error: Incidencia no encontrada")
		pausar()
		return
	}

	fmt.Println("\n--- Mecánicos disponibles ---")
	mecanicosDisponibles := []*Mecanico{}
	for _, m := range mecanicos {
		if m.Activo && m.Especialidad == incidencia.Tipo {
			mecanicosDisponibles = append(mecanicosDisponibles, m)
			fmt.Printf("ID: %d - %s (%d años exp, %d incidencias)\n",
				m.ID, m.Nombre, m.AniosExp, len(m.Incidencias))
		}
	}

	if len(mecanicosDisponibles) == 0 {
		fmt.Println("No hay mecánicos disponibles con la especialidad requerida")
		pausar()
		return
	}

	var idMecanico int
	fmt.Print("\nID del mecánico a asignar: ")
	fmt.Scanf("%d", &idMecanico)
	fmt.Scanln()

	var mecanico *Mecanico
	for _, m := range mecanicosDisponibles {
		if m.ID == idMecanico {
			mecanico = m
			break
		}
	}

	if mecanico == nil {
		fmt.Println("Error: Mecánico no encontrado o no disponible")
		pausar()
		return
	}

	// Verificar que no esté ya asignado
	for _, m := range incidencia.Mecanicos {
		if m.ID == mecanico.ID {
			fmt.Println("Error: El mecánico ya está asignado a esta incidencia")
			pausar()
			return
		}
	}

	incidencia.Mecanicos = append(incidencia.Mecanicos, mecanico)
	mecanico.Incidencias = append(mecanico.Incidencias, incidencia)

	fmt.Println("\nMecánico asignado exitosamente")
	pausar()
}

// Funciones de listado

func listarIncidenciasVehiculo() {
	limpiarPantalla()
	fmt.Println("=== INCIDENCIAS DE UN VEHÍCULO ===")

	var matricula string
	fmt.Print("Matrícula del vehículo: ")
	fmt.Scanln(&matricula)

	var vehiculo *Vehiculo
	for _, v := range vehiculos {
		if v.Matricula == matricula {
			vehiculo = v
			break
		}
	}

	if vehiculo == nil {
		fmt.Println("Error: Vehículo no encontrado")
		pausar()
		return
	}

	fmt.Printf("\nVehículo: %s %s (Matrícula: %s)\n",
		vehiculo.Marca, vehiculo.Modelo, vehiculo.Matricula)

	if vehiculo.Incidencia == nil {
		fmt.Println("Este vehículo no tiene incidencias")
	} else {
		inc := vehiculo.Incidencia
		fmt.Printf("\nID: %d\n", inc.ID)
		fmt.Printf("Tipo: %s\n", inc.Tipo)
		fmt.Printf("Prioridad: %s\n", inc.Prioridad)
		fmt.Printf("Estado: %s\n", inc.Estado)
		fmt.Printf("Descripción: %s\n", inc.Descripcion)

		if len(inc.Mecanicos) > 0 {
			fmt.Print("Mecánicos asignados: ")
			for i, m := range inc.Mecanicos {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(m.Nombre)
			}
			fmt.Println()
		}
	}

	pausar()
}

func listarVehiculosCliente() {
	limpiarPantalla()
	fmt.Println("=== VEHÍCULOS DE UN CLIENTE ===")

	var id int
	fmt.Print("ID del cliente: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	var cliente *Cliente
	for _, c := range clientes {
		if c.ID == id {
			cliente = c
			break
		}
	}

	if cliente == nil {
		fmt.Println("Error: Cliente no encontrado")
		pausar()
		return
	}

	fmt.Printf("\nCliente: %s\n", cliente.Nombre)
	fmt.Printf("Teléfono: %s\n", cliente.Telefono)
	fmt.Printf("Email: %s\n", cliente.Email)

	if cliente.Vehiculo == nil {
		fmt.Println("\nEste cliente no tiene vehículos registrados")
	} else {
		v := cliente.Vehiculo
		fmt.Printf("\n--- Vehículo ---\n")
		fmt.Printf("Matrícula: %s\n", v.Matricula)
		fmt.Printf("Marca: %s\n", v.Marca)
		fmt.Printf("Modelo: %s\n", v.Modelo)
		if v.EnTaller {
			fmt.Printf("Estado: En taller (Plaza %d)\n", v.NumeroPlaza)
		} else {
			fmt.Println("Estado: Fuera del taller")
		}
	}

	pausar()
}

func listarMecanicosDisponibles() {
	limpiarPantalla()
	fmt.Println("=== MECÁNICOS DISPONIBLES ===")

	disponibles := false
	for _, m := range mecanicos {
		if m.Activo && len(m.Incidencias) == 0 {
			disponibles = true
			fmt.Printf("\nID: %d\n", m.ID)
			fmt.Printf("Nombre: %s\n", m.Nombre)
			fmt.Printf("Especialidad: %s\n", m.Especialidad)
			fmt.Printf("Años de experiencia: %d\n", m.AniosExp)
			fmt.Println("---")
		}
	}

	if !disponibles {
		fmt.Println("No hay mecánicos disponibles sin incidencias asignadas")
	}

	pausar()
}

func listarIncidenciasMecanico() {
	limpiarPantalla()
	fmt.Println("=== INCIDENCIAS DE UN MECÁNICO ===")

	var id int
	fmt.Print("ID del mecánico: ")
	fmt.Scanf("%d", &id)
	fmt.Scanln()

	var mecanico *Mecanico
	for _, m := range mecanicos {
		if m.ID == id {
			mecanico = m
			break
		}
	}

	if mecanico == nil {
		fmt.Println("Error: Mecánico no encontrado")
		pausar()
		return
	}

	fmt.Printf("\nMecánico: %s (%s)\n", mecanico.Nombre, mecanico.Especialidad)

	if len(mecanico.Incidencias) == 0 {
		fmt.Println("Este mecánico no tiene incidencias asignadas")
	} else {
		fmt.Printf("\nTotal de incidencias: %d\n", len(mecanico.Incidencias))
		for _, inc := range mecanico.Incidencias {
			fmt.Printf("\n--- Incidencia ID: %d ---\n", inc.ID)
			fmt.Printf("Tipo: %s\n", inc.Tipo)
			fmt.Printf("Prioridad: %s\n", inc.Prioridad)
			fmt.Printf("Estado: %s\n", inc.Estado)
			fmt.Printf("Descripción: %s\n", inc.Descripcion)
		}
	}

	pausar()
}

func listarClientesConVehiculosEnTaller() {
	limpiarPantalla()
	fmt.Println("=== CLIENTES CON VEHÍCULOS EN TALLER ===")

	encontrados := false
	for _, c := range clientes {
		if c.Vehiculo != nil && c.Vehiculo.EnTaller {
			encontrados = true
			fmt.Printf("\n--- Cliente ---\n")
			fmt.Printf("ID: %d\n", c.ID)
			fmt.Printf("Nombre: %s\n", c.Nombre)
			fmt.Printf("Teléfono: %s\n", c.Telefono)
			fmt.Printf("Email: %s\n", c.Email)

			v := c.Vehiculo
			fmt.Printf("\nVehículo: %s %s\n", v.Marca, v.Modelo)
			fmt.Printf("Matrícula: %s\n", v.Matricula)
			fmt.Printf("Plaza: %d\n", v.NumeroPlaza)

			if v.Incidencia != nil {
				fmt.Printf("Incidencia: %s (%s)\n",
					v.Incidencia.Tipo, v.Incidencia.Estado)
			}
			fmt.Println("---")
		}
	}

	if !encontrados {
		fmt.Println("No hay clientes con vehículos en el taller actualmente")
	}

	pausar()
}

func listarTodasIncidenciasTaller() {
	limpiarPantalla()
	fmt.Println("=== TODAS LAS INCIDENCIAS DEL TALLER ===")

	if len(incidencias) == 0 {
		fmt.Println("No hay incidencias registradas en el taller")
		pausar()
		return
	}

	// Agrupar por estado
	abiertas := []*Incidencia{}
	enProceso := []*Incidencia{}
	cerradas := []*Incidencia{}

	for _, inc := range incidencias {
		switch inc.Estado {
		case "abierta":
			abiertas = append(abiertas, inc)
		case "en proceso":
			enProceso = append(enProceso, inc)
		case "cerrada":
			cerradas = append(cerradas, inc)
		}
	}

	fmt.Println("\n=== INCIDENCIAS ABIERTAS ===")
	if len(abiertas) == 0 {
		fmt.Println("Ninguna")
	} else {
		for _, inc := range abiertas {
			fmt.Printf("ID: %d - %s (%s) - %s\n",
				inc.ID, inc.Tipo, inc.Prioridad, inc.Descripcion)
		}
	}

	fmt.Println("\n=== INCIDENCIAS EN PROCESO ===")
	if len(enProceso) == 0 {
		fmt.Println("Ninguna")
	} else {
		for _, inc := range enProceso {
			fmt.Printf("ID: %d - %s (%s) - %s\n",
				inc.ID, inc.Tipo, inc.Prioridad, inc.Descripcion)
		}
	}

	fmt.Println("\n=== INCIDENCIAS CERRADAS ===")
	if len(cerradas) == 0 {
		fmt.Println("Ninguna")
	} else {
		for _, inc := range cerradas {
			fmt.Printf("ID: %d - %s (%s) - %s\n",
				inc.ID, inc.Tipo, inc.Prioridad, inc.Descripcion)
		}
	}

	fmt.Printf("\nTotal: %d incidencias (%d abiertas, %d en proceso, %d cerradas)\n",
		len(incidencias), len(abiertas), len(enProceso), len(cerradas))

	pausar()
}

// Menús

func menuClientes() {
	for {
		limpiarPantalla()
		fmt.Println("=== GESTIÓN DE CLIENTES ===")
		fmt.Println("1. Crear cliente")
		fmt.Println("2. Visualizar clientes")
		fmt.Println("3. Modificar cliente")
		fmt.Println("4. Eliminar cliente")
		fmt.Println("5. Listar vehículos de un cliente")
		fmt.Println("0. Volver al menú principal")

		var opcion int
		fmt.Print("\nSeleccione una opción: ")
		fmt.Scanf("%d", &opcion)
		fmt.Scanln()

		switch opcion {
		case 1:
			crearCliente()
		case 2:
			visualizarClientes()
		case 3:
			modificarCliente()
		case 4:
			eliminarCliente()
		case 5:
			listarVehiculosCliente()
		case 0:
			return
		default:
			fmt.Println("Opción inválida")
			pausar()
		}
	}
}

func menuVehiculos() {
	for {
		limpiarPantalla()
		fmt.Println("=== GESTIÓN DE VEHÍCULOS ===")
		fmt.Println("1. Crear vehículo")
		fmt.Println("2. Visualizar vehículos")
		fmt.Println("3. Modificar vehículo")
		fmt.Println("4. Eliminar vehículo")
		fmt.Println("5. Listar incidencias de un vehículo")
		fmt.Println("0. Volver al menú principal")

		var opcion int
		fmt.Print("\nSeleccione una opción: ")
		fmt.Scanf("%d", &opcion)
		fmt.Scanln()

		switch opcion {
		case 1:
			crearVehiculo()
		case 2:
			visualizarVehiculos()
		case 3:
			modificarVehiculo()
		case 4:
			eliminarVehiculo()
		case 5:
			listarIncidenciasVehiculo()
		case 0:
			return
		default:
			fmt.Println("Opción inválida")
			pausar()
		}
	}
}

func menuIncidencias() {
	for {
		limpiarPantalla()
		fmt.Println("=== GESTIÓN DE INCIDENCIAS ===")
		fmt.Println("1. Crear incidencia")
		fmt.Println("2. Visualizar incidencias")
		fmt.Println("3. Modificar incidencia")
		fmt.Println("4. Eliminar incidencia")
		fmt.Println("5. Cambiar estado de incidencia")
		fmt.Println("6. Asignar mecánico a incidencia")
		fmt.Println("7. Listar todas las incidencias del taller")
		fmt.Println("0. Volver al menú principal")

		var opcion int
		fmt.Print("\nSeleccione una opción: ")
		fmt.Scanf("%d", &opcion)
		fmt.Scanln()

		switch opcion {
		case 1:
			crearIncidencia()
		case 2:
			visualizarIncidencias()
		case 3:
			modificarIncidencia()
		case 4:
			eliminarIncidencia()
		case 5:
			cambiarEstadoIncidencia()
		case 6:
			asignarMecanicoAIncidencia()
		case 7:
			listarTodasIncidenciasTaller()
		case 0:
			return
		default:
			fmt.Println("Opción inválida")
			pausar()
		}
	}
}

func menuMecanicos() {
	for {
		limpiarPantalla()
		fmt.Println("=== GESTIÓN DE MECÁNICOS ===")
		fmt.Println("1. Crear mecánico")
		fmt.Println("2. Visualizar mecánicos")
		fmt.Println("3. Modificar mecánico")
		fmt.Println("4. Eliminar mecánico")
		fmt.Println("5. Dar alta/baja a mecánico")
		fmt.Println("6. Listar mecánicos disponibles")
		fmt.Println("7. Listar incidencias de un mecánico")
		fmt.Println("0. Volver al menú principal")

		var opcion int
		fmt.Print("\nSeleccione una opción: ")
		fmt.Scanf("%d", &opcion)
		fmt.Scanln()

		switch opcion {
		case 1:
			crearMecanico()
		case 2:
			visualizarMecanicos()
		case 3:
			modificarMecanico()
		case 4:
			eliminarMecanico()
		case 5:
			darAltaBajaMecanico()
		case 6:
			listarMecanicosDisponibles()
		case 7:
			listarIncidenciasMecanico()
		case 0:
			return
		default:
			fmt.Println("Opción inválida")
			pausar()
		}
	}
}

func menuTaller() {
	for {
		limpiarPantalla()
		fmt.Println("=== GESTIÓN DEL TALLER ===")
		fmt.Println("1. Asignar vehículo a taller")
		fmt.Println("2. Visualizar estado del taller")
		fmt.Println("3. Listar clientes con vehículos en taller")
		fmt.Println("0. Volver al menú principal")

		var opcion int
		fmt.Print("\nSeleccione una opción: ")
		fmt.Scanf("%d", &opcion)
		fmt.Scanln()

		switch opcion {
		case 1:
			asignarVehiculoATaller()
		case 2:
			visualizarEstadoTaller()
		case 3:
			listarClientesConVehiculosEnTaller()
		case 0:
			return
		default:
			fmt.Println("Opción inválida")
			pausar()
		}
	}
}

// *******************************************************************************
// Datos de prueba (opcional)
// *******************************************************************************

func cargarDatosPrueba() {
	// Crear mecánicos
	mec1 := &Mecanico{
		ID:           contadorMecanico,
		Nombre:       "Juan Pérez",
		Especialidad: "mecánica",
		AniosExp:     10,
		Activo:       true,
		Incidencias:  []*Incidencia{},
	}
	mecanicos = append(mecanicos, mec1)
	taller.Mecanicos = append(taller.Mecanicos, mec1)
	contadorMecanico++

	mec2 := &Mecanico{
		ID:           contadorMecanico,
		Nombre:       "María García",
		Especialidad: "eléctrica",
		AniosExp:     8,
		Activo:       true,
		Incidencias:  []*Incidencia{},
	}
	mecanicos = append(mecanicos, mec2)
	taller.Mecanicos = append(taller.Mecanicos, mec2)
	contadorMecanico++

	mec3 := &Mecanico{
		ID:           contadorMecanico,
		Nombre:       "Carlos López",
		Especialidad: "carrocería",
		AniosExp:     5,
		Activo:       true,
		Incidencias:  []*Incidencia{},
	}
	mecanicos = append(mecanicos, mec3)
	taller.Mecanicos = append(taller.Mecanicos, mec3)
	contadorMecanico++

	// Crear clientes
	cliente1 := &Cliente{
		ID:       contadorCliente,
		Nombre:   "Ana Martínez",
		Telefono: "600111222",
		Email:    "ana@email.com",
		Vehiculo: nil,
	}
	clientes = append(clientes, cliente1)
	contadorCliente++

	cliente2 := &Cliente{
		ID:       contadorCliente,
		Nombre:   "Pedro Sánchez",
		Telefono: "600333444",
		Email:    "pedro@email.com",
		Vehiculo: nil,
	}
	clientes = append(clientes, cliente2)
	contadorCliente++

	// Crear vehículos
	vehiculo1 := &Vehiculo{
		Matricula:    "1234ABC",
		Marca:        "Seat",
		Modelo:       "León",
		FechaEntrada: obtenerFechaActual(),
		FechaSalida:  "",
		Incidencia:   nil,
		EnTaller:     false,
		NumeroPlaza:  -1,
	}
	vehiculos = append(vehiculos, vehiculo1)
	cliente1.Vehiculo = vehiculo1

	vehiculo2 := &Vehiculo{
		Matricula:    "5678XYZ",
		Marca:        "Volkswagen",
		Modelo:       "Golf",
		FechaEntrada: obtenerFechaActual(),
		FechaSalida:  "",
		Incidencia:   nil,
		EnTaller:     false,
		NumeroPlaza:  -1,
	}
	vehiculos = append(vehiculos, vehiculo2)
	cliente2.Vehiculo = vehiculo2

	// Crear incidencias
	inc1 := &Incidencia{
		ID:          contadorIncidencia,
		Mecanicos:   []*Mecanico{},
		Tipo:        "mecánica",
		Prioridad:   "alta",
		Descripcion: "Cambio de correa de distribución",
		Estado:      "abierta",
	}
	incidencias = append(incidencias, inc1)
	vehiculo1.Incidencia = inc1
	contadorIncidencia++

	taller.TotalPlazas = calcularTotalPlazas()

	fmt.Println("Datos de prueba cargados exitosamente")
	pausar()
}

// *******************************************************************************

func inicializarSistema() {
	taller = Taller{
		Mecanicos:         []*Mecanico{},
		PlazasPorMecanico: 2,
		PlazasOcupadas:    make(map[int]bool),
		TotalPlazas:       0,
	}
}

func main() {
	inicializarSistema()

	for {
		limpiarPantalla()
		fmt.Println("╔════════════════════════════════════════╗")
		fmt.Println("║										║")
		fmt.Println("║   	SISTEMA DE GESTION DE TALLER     ║")
		fmt.Println("║          	MECANICO                 ║")
		fmt.Println("║										║")
		fmt.Println("╚════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("1. Gestión de Clientes")
		fmt.Println("2. Gestión de Vehículos")
		fmt.Println("3. Gestión de Incidencias")
		fmt.Println("4. Gestión de Mecánicos")
		fmt.Println("5. Gestión del Taller")
		fmt.Println("6. Cargar datos de prueba")
		fmt.Println("0. Salir")

		var opcion int
		fmt.Print("\nSeleccione una opción: ")
		fmt.Scanf("%d", &opcion)
		fmt.Scanln()

		switch opcion {
		case 1:
			menuClientes()
		case 2:
			menuVehiculos()
		case 3:
			menuIncidencias()
		case 4:
			menuMecanicos()
		case 5:
			menuTaller()
		case 6:
			cargarDatosPrueba()
		case 0:
			limpiarPantalla()
			fmt.Println("Gracias por usar el sistema. ¡Hasta pronto!")
			return
		default:
			fmt.Println("Opción inválida")
			pausar()
		}
	}
}
