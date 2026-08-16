// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   main.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/14 19:55:24 by dande-je          #+#    #+#             //
//   Updated: 2026/08/16 10:52:25 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

import (
	"fmt"
	"os"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/config"
)

func main() {
	if err := run(os.Args[0], os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "spider: %v\n", err)
		os.Exit(1)
	}
}

func run(progName string, args []string) error {
	_, err := config.Parse(progName, args)
	if err != nil {
		return err
	}

	return nil
}
