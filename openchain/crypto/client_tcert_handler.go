package crypto

import (
	"crypto/x509"
	"github.com/openblockchain/obc-peer/openchain/crypto/utils"
	obc "github.com/openblockchain/obc-peer/protos"
)

type tCertHandlerImpl struct {
	client    *clientImpl
	txHandler *tCertTransactionHandlerImpl

	tCert *x509.Certificate
}

type tCertTransactionHandlerImpl struct {
	tCertHandler *tCertHandlerImpl
	hook         []byte
}

func (handler *tCertHandlerImpl) initDER(client *clientImpl, tCertDER []byte) error {
	// Parse the DER
	tCert, err := utils.DERToX509Certificate(tCertDER)
	if err != nil {
		client.node.log.Error("Failed parsing TCert DER [%s].", err.Error())

		return err
	}

	handler.initX509(client, tCert)

	return nil
}

func (handler *tCertHandlerImpl) initX509(client *clientImpl, tCert *x509.Certificate) error {
	handler.client = client
	handler.tCert = tCert
	handler.txHandler = &tCertTransactionHandlerImpl{}
	handler.txHandler.init(handler)

	return nil
}

// GetCertificate returns the TCert DER
func (handler *tCertHandlerImpl) GetCertificate() []byte {
	return utils.Clone(handler.tCert.Raw)
}

// Sign signs msg using the signing key corresponding to this TCert
func (handler *tCertHandlerImpl) Sign(msg []byte) ([]byte, error) {
	return handler.client.signUsingTCertX509(handler.tCert, msg)
}

// Verify verifies msg using the verifying key corresponding to this TCert
func (handler *tCertHandlerImpl) Verify(signature []byte, msg []byte) error {
	return handler.client.verifyUsingTCertX509(handler.tCert, signature, msg)
}

// GetTransactionHandler returns the transaction handler relative to this certificate
func (handler *tCertHandlerImpl) GetTransactionHandler() TransactionHandler {
	// TODO: in the future should return a different instance
	return handler.txHandler
}

func (handler *tCertTransactionHandlerImpl) init(tCertHandler *tCertHandlerImpl) error {
	handler.tCertHandler = tCertHandler
	handler.hook = utils.Hash(handler.tCertHandler.tCert.Raw)

	return nil
}

// GetCertificateHandler returns the certificate handler relative to the certificate mapped to this transaction
func (handler *tCertTransactionHandlerImpl) GetCertificateHandler() (CertificateHandler, error) {
	return handler.tCertHandler, nil
}

// GetHook returns an Hook to the underlying transaction layer
func (handler *tCertTransactionHandlerImpl) GetHook() ([]byte, error) {
	return utils.Clone(handler.hook), nil
}

// NewChaincodeDeployTransaction is used to deploy chaincode.
func (handler *tCertTransactionHandlerImpl) NewChaincodeDeployTransaction(chaincodeDeploymentSpec *obc.ChaincodeDeploymentSpec, uuid string) (*obc.Transaction, error) {
	return handler.tCertHandler.client.newChaincodeDeployUsingTCert(chaincodeDeploymentSpec, uuid, handler.tCertHandler.tCert.Raw)
}

// NewChaincodeExecute is used to execute chaincode's functions.
func (handler *tCertTransactionHandlerImpl) NewChaincodeExecute(chaincodeInvocation *obc.ChaincodeInvocationSpec, uuid string) (*obc.Transaction, error) {
	return handler.tCertHandler.client.newChaincodeExecuteUsingTCert(chaincodeInvocation, uuid, handler.tCertHandler.tCert.Raw)
}

// NewChaincodeQuery is used to query chaincode's functions.
func (handler *tCertTransactionHandlerImpl) NewChaincodeQuery(chaincodeInvocation *obc.ChaincodeInvocationSpec, uuid string) (*obc.Transaction, error) {
	return handler.tCertHandler.client.newChaincodeQueryUsingTCert(chaincodeInvocation, uuid, handler.tCertHandler.tCert.Raw)
}
