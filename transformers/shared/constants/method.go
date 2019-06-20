// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package constants

//TODO: get cat and jug file method signatures directly from the ABI
func biteMethod() string               { return getSolidityFunctionSignature(CatABI(), "Bite") }
func catFileChopLumpMethod() string    { return "file(bytes32,bytes32,uint256)" }
func catFileFlipMethod() string        { return getSolidityFunctionSignature(CatABI(), "file") }
func catFileVowMethod() string         { return "file(bytes32,address)" }
func dealMethod() string               { return getSolidityFunctionSignature(FlipperABI(), "deal") }
func dentMethod() string               { return getSolidityFunctionSignature(FlipperABI(), "dent") }
func flapKickMethod() string           { return getSolidityFunctionSignature(FlapperABI(), "Kick") }
func flipKickMethod() string           { return getSolidityFunctionSignature(FlipperABI(), "Kick") }
func flopKickMethod() string           { return getSolidityFunctionSignature(FlopperABI(), "Kick") }
func jugDripMethod() string            { return getSolidityFunctionSignature(JugABI(), "drip") }
func jugFileBaseMethod() string        { return "file(bytes32,uint256)" }
func jugFileIlkMethod() string         { return "file(bytes32,bytes32,uint256)" }
func jugFileVowMethod() string         { return "file(bytes32,address)" }
func jugInitMethod() string            { return getSolidityFunctionSignature(JugABI(), "init") }
func spotFileMatMethod() string        { return "file(bytes32,bytes32,uint256)" }
func spotFilePipMethod() string        { return "file(bytes32,address)" }
func spotPokeMethod() string           { return getSolidityFunctionSignature(SpotABI(), "Poke") }
func tendMethod() string               { return getSolidityFunctionSignature(FlipperABI(), "tend") }
func vatFileDebtCeilingMethod() string { return "file(bytes32,uint256)" }
func vatFileIlkMethod() string         { return "file(bytes32,bytes32,uint256)" }
func vatFluxMethod() string            { return getSolidityFunctionSignature(VatABI(), "flux") }
func vatFoldMethod() string            { return getSolidityFunctionSignature(VatABI(), "fold") }
func vatForkMethod() string            { return getSolidityFunctionSignature(VatABI(), "fork") }
func vatFrobMethod() string            { return getSolidityFunctionSignature(VatABI(), "frob") }
func vatGrabMethod() string            { return getSolidityFunctionSignature(VatABI(), "grab") }
func vatHealMethod() string            { return getSolidityFunctionSignature(VatABI(), "heal") }
func vatInitMethod() string            { return getSolidityFunctionSignature(VatABI(), "init") }
func vatMoveMethod() string            { return getSolidityFunctionSignature(VatABI(), "move") }
func vatSlipMethod() string            { return getSolidityFunctionSignature(VatABI(), "slip") }
func vatSuckMethod() string            { return getSolidityFunctionSignature(VatABI(), "suck") }
func vowFessMethod() string            { return getSolidityFunctionSignature(VowABI(), "fess") }
func vowFileMethod() string            { return getSolidityFunctionSignature(VowABI(), "file") }
func vowFlogMethod() string            { return getSolidityFunctionSignature(VowABI(), "flog") }
