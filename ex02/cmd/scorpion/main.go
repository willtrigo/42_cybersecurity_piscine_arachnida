// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   main.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/21 14:36:07 by dande-je          #+#    #+#             //
//   Updated: 2026/08/26 08:50:35 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

import (
	"fmt"
	"os"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/config"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/parser"
)

func main() {
	if err := run(os.Args[0], os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "scorpion: %v\n", err)
		os.Exit(1)
	}
}

func run(progName string, args []string) error {
	_, err := config.Parse(progName, args)
	if err != nil {
		return err
	}

	_ = parser.NewRegistry()

	return nil
}
